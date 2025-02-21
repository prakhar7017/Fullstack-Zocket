export interface User {
  id: number;
  name: string;
  email?: string;
}

export interface Task {
  id: number;
  title: string;
  description: string;
  assignee_id: number;
  status: 'PENDING' | 'IN_PROGRESS' | 'COMPLETED';
  created_at: string;
  updated_at: string;
  deadline: string;
  importance: number;
}

export interface AuthResponse {
  token: string;
}

export interface UserResponse {
  user: {
    id: number;
    name: string;
    email: string;
  };
}

export interface UsersResponse {
  users: User[];
}

export interface TaskEvent {
  event: string;
  task?: Task;
  tasks?: Task[];
}

export interface AIResponse {
  suggestions?: string[];
  breakdown?: string[];
  prioritized_tasks?: Task[];
}