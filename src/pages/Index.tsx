
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import EmailList from "@/components/EmailList";
import CategorySidebar from "@/components/CategorySidebar";
import LoginScreen from "@/components/LoginScreen";

const Index = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  if (!isAuthenticated) {
    return <LoginScreen onLogin={() => setIsAuthenticated(true)} />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100">
      <div className="container mx-auto px-4 py-8">
        <div className="flex gap-6">
          <CategorySidebar />
          <div className="flex-1">
            <EmailList />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Index;
