import { AlertCircle } from "lucide-react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

function Maintenance() {
  return (
    <div className="w-full flex justify-center items-center">
      <Alert variant="default" className="max-w-sm flex">
        <AlertCircle className="h-7 w-7 pr-2" />
        <div className="flex flex-col justify-center items-start">
          <AlertTitle className="text-2xl">Sorry</AlertTitle>
          <AlertDescription className="text-lg">
            We are under maintenance. Please try again later.
          </AlertDescription>
        </div>
      </Alert>
    </div>
  );
}

export default Maintenance;
