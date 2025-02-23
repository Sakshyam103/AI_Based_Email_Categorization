
import { Card } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";

const mockEmails = [
  {
    id: 1,
    subject: "CPT Application Update",
    sender: "iso@oswego.edu",
    preview: "Your CPT application has been processed...",
    category: "cpt",
    date: "10:30 AM",
  },
  {
    id: 2,
    subject: "OPT Document Requirements",
    sender: "advisors@oswego.edu",
    preview: "Please submit the following documents...",
    category: "opt",
    date: "Yesterday",
  },
];

const EmailList = () => {
  return (
    <Card className="h-[calc(100vh-4rem)] backdrop-blur-lg bg-white/90 shadow-lg">
      <ScrollArea className="h-full">
        <div className="p-4">
          <h2 className="text-lg font-semibold mb-4">Inbox</h2>
          <div className="space-y-2">
            {mockEmails.map((email) => (
              <Card
                key={email.id}
                className="p-4 hover:bg-gray-50 cursor-pointer transition-all duration-200 animate-fade-in"
              >
                <div className="flex items-start justify-between">
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <h3 className="font-medium text-gray-900">
                        {email.subject}
                      </h3>
                      <Badge
                        variant="secondary"
                        className={`bg-category-${email.category} text-white`}
                      >
                        {email.category.toUpperCase()}
                      </Badge>
                    </div>
                    <p className="text-sm text-gray-600">{email.sender}</p>
                    <p className="text-sm text-gray-500">{email.preview}</p>
                  </div>
                  <span className="text-xs text-gray-400">{email.date}</span>
                </div>
              </Card>
            ))}
          </div>
        </div>
      </ScrollArea>
    </Card>
  );
};

export default EmailList;
