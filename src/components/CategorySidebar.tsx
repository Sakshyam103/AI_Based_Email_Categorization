
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn } from "@/lib/utils";
import {
  FileText,
  Briefcase,
  KeyRound,
  CalendarCheck,
  FileSignature,
  Heart,
  Plane,
  Files,
  Award,
  Clock,
} from "lucide-react";

const categories = [
  { id: "cpt", name: "CPT", icon: Briefcase, color: "category-cpt" },
  { id: "opt", name: "OPT", icon: FileText, color: "category-opt" },
  { id: "ssn", name: "SSN", icon: KeyRound, color: "category-ssn" },
  {
    id: "registration",
    name: "Pre-Registration",
    icon: CalendarCheck,
    color: "category-registration",
  },
  {
    id: "application",
    name: "Application",
    icon: FileSignature,
    color: "category-application",
  },
  {
    id: "health",
    name: "Health Insurance",
    icon: Heart,
    color: "category-health",
  },
  { id: "arrival", name: "Arrival", icon: Plane, color: "category-arrival" },
  {
    id: "documents",
    name: "Documents",
    icon: Files,
    color: "category-documents",
  },
  { id: "stemopt", name: "STEM OPT", icon: Award, color: "category-stemopt" },
  {
    id: "stemext",
    name: "STEM OPT Extension",
    icon: Clock,
    color: "category-stemopt",
  },
];

const CategorySidebar = () => {
  return (
    <Card className="w-64 h-[calc(100vh-4rem)] backdrop-blur-lg bg-white/90 shadow-lg">
      <ScrollArea className="h-full">
        <div className="p-4 space-y-2">
          <h2 className="text-lg font-semibold mb-4">Categories</h2>
          {categories.map((category) => (
            <Button
              key={category.id}
              variant="ghost"
              className={cn(
                "w-full justify-start gap-2 hover:bg-gray-100 transition-all duration-200",
                "text-gray-700 font-medium"
              )}
            >
              <category.icon
                className={cn("w-4 h-4", `text-${category.color}`)}
              />
              {category.name}
            </Button>
          ))}
        </div>
      </ScrollArea>
    </Card>
  );
};

export default CategorySidebar;
