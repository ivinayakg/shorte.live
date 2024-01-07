import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Link } from "react-router-dom";
import { useMain } from "@/components/main-provider";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import fetch from "@/utils/axios";
import { useToast } from "@/components/ui/use-toast";

function CreateShort() {
  const { userState } = useMain();
  const loginWithGoogleUrl = import.meta.env.VITE_GOOGLE_SIGN_IN;
  const { toast } = useToast();

  const generateShortLink = async (e: any) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    try {
      const requestData = {
        destination: formData.get("destination"),
        short: formData.get("short"),
        expiry: 96,
      };

      const res = await fetch.post("/url", requestData, {
        headers: { Authorization: `Bearer ${userState.token}` },
      });
      if (res.status !== 201) {
        throw new Error("Check again later");
      }
      const data = res.data;
      navigator.clipboard.writeText(data.short);
      toast({
        title: "Short URL generated Successfully",
        description: "The URL is copied to your clipboard",
        duration: 2000,
      });
    } catch (error: any) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Uh oh! Something went wrong.",
        description: <p>{error.message}</p>,
        duration: 2000,
      });
    }
  };

  return (
    <div className="w-full flex flex-col justify-center items-center gap-10 createshort">
      <h2 className="text-7xl font-bold">Generate shorten link</h2>
      <form
        className="flex w-full max-w-5xl items-center space-x-2"
        onSubmit={generateShortLink}
      >
        <Input
          className=""
          type="text"
          placeholder="Destination URL"
          name="destination"
        />
        <Input
          className="w-2/5"
          type="text"
          placeholder="Custom Short (Optional)"
          name="short"
        />
        {userState.login ? (
          <Button className="" type="submit">
            Generate
          </Button>
        ) : (
          <>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button>Generate</Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Login First?</AlertDialogTitle>
                  <AlertDialogDescription>
                    Our service is free, but only to verified users!
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction>
                    <Link to={loginWithGoogleUrl}>
                      <Button className="gap-2">
                        Login With Google <img src="/google.svg" alt="" />
                      </Button>
                    </Link>
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </>
        )}
      </form>
    </div>
  );
}

export default CreateShort;
