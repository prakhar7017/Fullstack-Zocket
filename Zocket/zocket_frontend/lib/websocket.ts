import { Task, TaskEvent } from './types';

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private token: string;
  private messageHandlers: ((event: TaskEvent) => void)[] = [];
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;

  constructor(token: string) {
    this.token = token;
    this.connect();
  }

  private connect() {
    try {
      this.ws = new WebSocket(`ws://localhost:8080/ws/?token=${this.token}`);
      
      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
      };
      
      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data) as TaskEvent;
          this.messageHandlers.forEach(handler => handler(data));
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      this.ws.onclose = () => {
        console.log('WebSocket disconnected');
        this.handleReconnect();
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      this.handleReconnect();
    }
  }

  private handleReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
      setTimeout(() => this.connect(), this.reconnectDelay * this.reconnectAttempts);
    } else {
      console.error('Max reconnection attempts reached');
    }
  }

  public onMessage(handler: (event: TaskEvent) => void) {
    this.messageHandlers.push(handler);
  }

  public createTask(task: Partial<Task>) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({
        action: 'create_task',
        ...task
      }));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  public updateTask(task: Partial<Task>) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({
        action: 'update_task',
        ...task
      }));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  public getTasks() {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({
        action: 'get_tasks'
      }));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  public close() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}