import { AlertCircle } from "lucide-react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

function NotFound() {
  return (
    <div className="w-full flex justify-center items-center">
      <Alert variant="destructive" className="max-w-sm flex">
        <AlertCircle className="h-7 w-7 pr-2" />
        <div className="flex flex-col justify-center items-start">
          <AlertTitle className="text-2xl">Error</AlertTitle>
          <AlertDescription className="text-lg">
            This link is either expired or invalid.
          </AlertDescription>
        </div>
      </Alert>
    </div>
  );
}

export default NotFound;
