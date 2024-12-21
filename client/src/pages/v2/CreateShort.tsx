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
import { HeadingOne, HeadingTwo, Paragraph } from "@/components/typography";
import { useState } from "react";

function CreateShort() {
  const [tempShortCompleted, setTempShortCompleted] = useState(false);

  const { userState } = useMain();
  const loginWithGoogleUrl = import.meta.env.VITE_GOOGLE_SIGN_IN;
  const { toast } = useToast();

  const generateShortLink = async (e: any) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    try {
      const inputTime = formData.get("expiry") as string;
      const inputTimeUnix = (new Date(inputTime).getTime() / 1000).toFixed();
      const requestData = {
        destination: formData.get("destination"),
        short: formData.get("short"),
        expiry: Number(inputTimeUnix),
      };

      const res = await fetch.post("/url", requestData, {
        withCredentials: true,
      });
      if (res.status !== 201) {
        throw new Error("Check again later");
      }
      const data = res.data;
      navigator.clipboard.writeText(data.short);
      toast({
        title: "Short URL generated Successfully",
        description: `The URL is copied to your clipboard, link -> ${data.short}`,
        duration: 2000,
      });

      e.target.reset();
    } catch (error: any) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Uh oh! Something went wrong.",
        description: <p>{error.response.data.error}</p>,
        duration: 2000,
      });
    }
  };

  const generateShortLinkTemp = async (e: any) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    try {
      const requestData = {
        destination: formData.get("destination"),
      };

      const res = await fetch.post("/url/temp", requestData, {
        withCredentials: true,
      });
      if (res.status !== 201) {
        if (res.status === 429) {
          setTempShortCompleted(true);
        }
        throw new Error("Check again later");
      }
      const data = res.data;
      navigator.clipboard.writeText(data.short);
      toast({
        title: "Short URL generated Successfully",
        description: `The URL is copied to your clipboard, link -> ${data.short}`,
        duration: 2000,
      });

      e.target.reset();
    } catch (error: any) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Uh oh! Something went wrong.",
        description: <p>{error.response.data.error}</p>,
        duration: 2000,
      });
    }
  };

  return (
    <div className="w-full flex flex-col justify-start items-center gap-2 sm:gap-5 createshort p-8">
      <div className="mt-36 flex flex-col justify-start items-center gap-2 sm:gap-5">
        <HeadingOne className="border-none lg:text-7xl sm:text-5xl font-semibold text-primary">
          Shorten your{" "}
          <span className="underline decoration-[#44C67F] underline-offset-8">
            URL
          </span>{" "}
          with just a Click
        </HeadingOne>
        <Paragraph className="text-xl text-[#8B8787] !mt-8">
          shorte.live, by Vinayak Goyal, makes URL shortening a breeze. Built
          with React.js and{" "}
          <span className="font-semibold text-primary">Golang</span>.
        </Paragraph>
      </div>
      <form
        className="flex flex-col w-full max-w-5xl items-center gap-6 justify-center mt-24"
        onSubmit={userState.login ? generateShortLink : generateShortLinkTemp}
      >
        {!userState.login && (
          <>
            <Paragraph className="text-xl font-semibold text-primary">
              Try it now!
            </Paragraph>
          </>
        )}
        <div className="flex w-full items-center space-x-2 flex-col gap-2 justify-center sm:flex-row ">
          <Input
            className={`w-80 ${userState.login && "sm:w-full"}`}
            type="text"
            placeholder="Destination URL"
            name="destination"
          />
          {userState.login && (
            <>
              <Input
                className="w-full sm:w-2/6"
                type="text"
                placeholder="Custom Short (Optional)"
                name="short"
                style={{ margin: "0" }}
                disabled={!userState.login}
              />
              <Input
                className="w-full sm:w-2/6"
                type="datetime-local"
                name="expiry"
                style={{ margin: "0" }}
                disabled={!userState.login}
              />
            </>
          )}
          {userState.login || !tempShortCompleted ? (
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
        </div>
      </form>
    </div>
  );
}

export default CreateShort;
