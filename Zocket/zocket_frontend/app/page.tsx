import { Button } from '@/components/ui/button';
import Link from 'next/link';
import { ArrowRight } from 'lucide-react';

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-primary/10 to-secondary/10 p-4">
      <div className="max-w-3xl text-center space-y-6">
        <h1 className="text-4xl sm:text-6xl font-bold text-primary">
          Task Management System
        </h1>
        <p className="text-lg sm:text-xl text-muted-foreground">
          Manage your tasks efficiently with real-time updates and AI-powered suggestions
        </p>
        <Link href="/auth">
          <Button size="lg" className="mt-6">
            Get Started <ArrowRight className="ml-2" />
          </Button>
        </Link>
      </div>
    </div>
  );
}