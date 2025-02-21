'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { WebSocketClient } from '@/lib/websocket';
import { Task, User } from '@/lib/types';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { getCurrentUser, getUsers, getSuggestions, getBreakdown, prioritizeTasks } from '@/lib/api';
import { Badge } from '@/components/ui/badge';
import { Calendar } from '@/components/ui/calendar';
import { Slider } from '@/components/ui/slider';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Brain, Calendar as CalendarIcon, ListTodo, Plus, RefreshCcw } from 'lucide-react';

export default function Dashboard() {
  const router = useRouter();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [wsClient, setWsClient] = useState<WebSocketClient | null>(null);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const [breakdown, setBreakdown] = useState<string[]>([]);
  const [prompt, setPrompt] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/auth');
      return;
    }

    const initializeData = async () => {
      try {
        const userResponse = await getCurrentUser(token);
        setCurrentUser(userResponse.user);

        const usersResponse = await getUsers(token);
        setUsers(usersResponse.users);

        const ws = new WebSocketClient(token);
        setWsClient(ws);

        ws.onMessage((event) => {
          if (event.event === 'task_list') {
            setTasks(event.tasks || []);
          } else if (event.event === 'task_created' || event.event === 'task_updated') {
            ws.getTasks();
          }
        });

        ws.getTasks();
      } catch (error) {
        console.error('Failed to initialize dashboard:', error);
        router.push('/auth');
      }
    };

    initializeData();

    return () => {
      if (wsClient) {
        wsClient.close();
      }
    };
  }, []);

  const handleCreateTask = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = e.currentTarget;
    const formData = new FormData(form);
    
    wsClient?.createTask({
      title: formData.get('title') as string,
      description: formData.get('description') as string,
      assignee_id: parseInt(formData.get('assignee_id') as string),
      importance: parseInt(formData.get('importance') as string),
      deadline: formData.get('deadline') as string,
      status: 'PENDING'
    });

    form.reset();
  };

  const handleUpdateTask = (task: Task, updates: Partial<Task>) => {
    wsClient?.updateTask({ ...task, ...updates });
  };

  const handleGetSuggestions = async () => {
    setIsLoading(true);
    try {
      const token = localStorage.getItem('token');
      if (!token) return;

      const response = await getSuggestions(prompt, token);
      setSuggestions(response.suggestions || []);
    } finally {
      setIsLoading(false);
    }
  };

  const handleGetBreakdown = async () => {
    setIsLoading(true);
    try {
      const token = localStorage.getItem('token');
      if (!token) return;

      const response = await getBreakdown(prompt, token);
      setBreakdown(response.breakdown || []);
    } finally {
      setIsLoading(false);
    }
  };

  const handlePrioritizeTasks = async () => {
    setIsLoading(true);
    try {
      const token = localStorage.getItem('token');
      if (!token) return;

      const response = await prioritizeTasks(tasks, token);
      if (response.prioritized_tasks) {
        setTasks(response.prioritized_tasks);
      }
    } finally {
      setIsLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'PENDING': return 'bg-yellow-500';
      case 'IN_PROGRESS': return 'bg-blue-500';
      case 'COMPLETED': return 'bg-green-500';
      default: return 'bg-gray-500';
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-secondary/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-center">
          <h1 className="text-4xl font-bold">Task Dashboard</h1>
          <div className="flex gap-2">
            <Dialog>
              <DialogTrigger asChild>
                <Button>
                  <Brain className="mr-2 h-4 w-4" />
                  AI Assistant
                </Button>
              </DialogTrigger>
              <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                  <DialogTitle>AI Assistant</DialogTitle>
                </DialogHeader>
                <div className="space-y-4">
                  <Textarea
                    placeholder="Enter your prompt..."
                    value={prompt}
                    onChange={(e) => setPrompt(e.target.value)}
                  />
                  <div className="flex gap-2">
                    <Button onClick={handleGetSuggestions} disabled={isLoading}>
                      Get Suggestions
                    </Button>
                    <Button onClick={handleGetBreakdown} disabled={isLoading}>
                      Get Breakdown
                    </Button>
                  </div>
                  {suggestions.length > 0 && (
                    <div>
                      <h3 className="font-semibold mb-2">Suggestions:</h3>
                      <ul className="list-disc pl-4">
                        {suggestions.map((suggestion, i) => (
                          <li key={i}>{suggestion}</li>
                        ))}
                      </ul>
                    </div>
                  )}
                  {breakdown.length > 0 && (
                    <div>
                      <h3 className="font-semibold mb-2">Breakdown:</h3>
                      <ul className="list-disc pl-4">
                        {breakdown.map((step, i) => (
                          <li key={i}>{step}</li>
                        ))}
                      </ul>
                    </div>
                  )}
                </div>
              </DialogContent>
            </Dialog>

            <Dialog>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  New Task
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create New Task</DialogTitle>
                </DialogHeader>
                <form onSubmit={handleCreateTask} className="space-y-4">
                  <div>
                    <Input name="title" placeholder="Task Title" required />
                  </div>
                  <div>
                    <Textarea name="description" placeholder="Task Description" required />
                  </div>
                  <div>
                    <Select name="assignee_id" required>
                      <SelectTrigger>
                        <SelectValue placeholder="Select Assignee" />
                      </SelectTrigger>
                      <SelectContent>
                        {users.map((user) => (
                          <SelectItem key={user.id} value={user.id.toString()}>
                            {user.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-2">Importance (1-5)</label>
                    <Slider
                      name="importance"
                      defaultValue={[3]}
                      min={1}
                      max={5}
                      step={1}
                      className="w-full"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-2">Deadline</label>
                    <Input type="date" name="deadline" required />
                  </div>
                  <Button type="submit" className="w-full">Create Task</Button>
                </form>
              </DialogContent>
            </Dialog>

            <Button onClick={handlePrioritizeTasks} disabled={isLoading}>
              <RefreshCcw className="mr-2 h-4 w-4" />
              Prioritize Tasks
            </Button>
          </div>
        </div>

        <ScrollArea className="h-[calc(100vh-12rem)]">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {tasks.map((task) => (
              <Card key={task.id} className="p-4 space-y-4">
                <div className="flex justify-between items-start">
                  <h3 className="text-lg font-semibold">{task.title}</h3>
                  <Badge className={getStatusColor(task.status)}>
                    {task.status}
                  </Badge>
                </div>
                <p className="text-muted-foreground">{task.description}</p>
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <CalendarIcon className="h-4 w-4" />
                  <span>Due: {new Date(task.deadline).toLocaleDateString()}</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-sm text-muted-foreground">
                    Importance: {task.importance}
                  </span>
                  <span className="text-sm text-muted-foreground">
                    Assignee: {users.find(u => u.id === task.assignee_id)?.name}
                  </span>
                </div>
                <div className="flex gap-2">
                  <Select
                    value={task.status}
                    onValueChange={(value) => handleUpdateTask(task, { status: value })}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="PENDING">Pending</SelectItem>
                      <SelectItem value="IN_PROGRESS">In Progress</SelectItem>
                      <SelectItem value="COMPLETED">Completed</SelectItem>
                    </SelectContent>
                  </Select>
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button variant="outline" size="sm">
                        <ListTodo className="h-4 w-4" />
                      </Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Task Actions</DialogTitle>
                      </DialogHeader>
                      <div className="space-y-4">
                        <Button
                          onClick={() => {
                            const token = localStorage.getItem('token');
                            if (!token) return;
                            getSuggestions(task.description, token)
                              .then(response => setSuggestions(response.suggestions || []));
                          }}
                          className="w-full"
                        >
                          Get Suggestions
                        </Button>
                        <Button
                          onClick={() => {
                            const token = localStorage.getItem('token');
                            if (!token) return;
                            getBreakdown(task.description, token)
                              .then(response => setBreakdown(response.breakdown || []));
                          }}
                          className="w-full"
                        >
                          Get Breakdown
                        </Button>
                      </div>
                    </DialogContent>
                  </Dialog>
                </div>
              </Card>
            ))}
          </div>
        </ScrollArea>
      </div>
    </div>
  );
}