
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";

interface LoginScreenProps {
  onLogin: () => void;
}

const LoginScreen = ({ onLogin }: LoginScreenProps) => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-gray-50 to-gray-100">
      <Card className="w-full max-w-md p-8 backdrop-blur-lg bg-white/90 shadow-lg animate-fade-in">
        <div className="text-center space-y-6">
          <h1 className="text-3xl font-semibold text-gray-900">
            Email Management System
          </h1>
          <p className="text-gray-600">
            Sign in with your Oswego email to continue
          </p>
          <Button
            onClick={onLogin}
            className="w-full bg-primary text-white hover:bg-primary/90 transition-all duration-200"
          >
            Sign in with Google
          </Button>
          <p className="text-sm text-gray-500">
            Only @oswego.edu accounts are allowed
          </p>
        </div>
      </Card>
    </div>
  );
};

export default LoginScreen;
