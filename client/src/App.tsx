import { Button } from "@/components/ui/button";
import { ThemeProvider } from "@/components/theme-provider";
import "./App.css";
import { ModeToggle } from "./components/mode-toggle";
import { Input } from "./components/ui/input";
import {
  createBrowserRouter,
  Outlet,
  RouterProvider,
  useNavigate,
  useSearchParams,
  Link,
  useLocation,
} from "react-router-dom";
import { MainProvider, useMain } from "./components/main-provider";
import { useEffect, useState } from "react";
import { setInLocalStorage } from "./utils/localstorage";
import { AvatarImage, Avatar, AvatarFallback } from "./components/ui/avatar";
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
import fetch from "./utils/axios";
import { Toaster } from "./components/ui/toaster";
import { useToast } from "./components/ui/use-toast";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Main />,
    children: [
      {
        path: "",
        element: <CreateShort />,
      },
      {
        path: "/my-urls",
        element: <MyUrls />,
      },
    ],
  },
  {
    path: "/auth/",
    element: <Auth />,
  },
]);

function Auth() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    const token = searchParams.get("token");
    if (token) {
      setInLocalStorage("userToken", token);
    }
    navigate("/");
  });

  return null;
}

function MyUrls() {
  const { userState } = useMain();
  const [urlsData, setUrlsData] = useState([]);
  const { toast } = useToast();
  const Headings = ["Serial", "Short", "Destination", "Expiry", "ID"];
  const location = useLocation();

  const CopyToClipboard = (value: string) => {
    navigator.clipboard.writeText(value);
    toast({
      title: "Successfully copied to your clipboard",
      duration: 2000,
    });
  };

  useEffect(() => {
    if (userState.login && location.pathname === "/my-urls") {
      (async () => {
        const res = await fetch.get("/url/all", {
          headers: { Authorization: `Bearer ${userState.token}` },
        });
        const data = res.data;
        setUrlsData(data);
        if (data.length === 0) {
          toast({
            title: "You have 0 links",
            duration: 2000,
          });
        }
      })();
    }
  }, [location.pathname, userState.login]);


  return (
    <div className="py-5">
      <Table>
        <TableHeader>
          <TableRow>
            {Headings.map((v, i) => {
              return (
                <TableHead className="w-[100px]" key={i}>
                  {v}
                </TableHead>
              );
            })}
          </TableRow>
        </TableHeader>
        <TableBody>
          {urlsData.length ? (
            urlsData.map((url: any, i) => {
              const date = new Date(url.expiry);
              return (
                <TableRow key={url._id}>
                  <TableCell className="font-medium text-left">
                    {i + 1}
                  </TableCell>
                  <TableCell
                    className="font-medium text-left hover:cursor-pointer hover:underline"
                    onClick={() => CopyToClipboard(url.short)}
                  >
                    {url.short}
                  </TableCell>
                  <TableCell
                    className="font-medium text-left hover:cursor-pointer hover:underline"
                    onClick={() => CopyToClipboard(url.destination)}
                  >
                    {url.destination}
                  </TableCell>
                  <TableCell className="font-medium text-left">
                    {date.toLocaleString()}
                  </TableCell>
                  <TableCell className="font-medium text-left">
                    {url._id}
                  </TableCell>
                </TableRow>
              );
            })
          ) : (
            <></>
          )}
        </TableBody>
      </Table>
    </div>
  );
}

function Header() {
  const loginWithGoogleUrl = import.meta.env.VITE_GOOGLE_SIGN_IN;
  const myUrls = "/my-urls";
  const { userState } = useMain();
  const navigate = useNavigate();

  const logout = () => {
    setInLocalStorage("userToken", null);
    navigate(0);
  };

  return (
    <div className="min-w-full flex justify-between items-center py-5">
      <Link to="/">
        <h1 className="text-3xl">URL Shortner</h1>
      </Link>
      <div className="flex justify-center items-center gap-4">
        {userState.login ? (
          <>
            <Avatar>
              <AvatarImage src={userState.picture} alt="@shadcn" />
              <AvatarFallback>U</AvatarFallback>
            </Avatar>
            <Button className="" onClick={logout}>
              Logout
            </Button>
            <Link to={myUrls}>
              <Button variant="ghost">My URLs</Button>
            </Link>
          </>
        ) : (
          <Link to={loginWithGoogleUrl}>
            <Button className="gap-2">
              Login With Google <img src="/google.svg" alt="" />
            </Button>
          </Link>
        )}
        <ModeToggle />
      </div>
    </div>
  );
}

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

function Main() {
  return (
    <>
      <Header />
      <Outlet />
    </>
  );
}

function App() {
  return (
    <>
      <MainProvider>
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
          <RouterProvider router={router} />
          <Toaster />
        </ThemeProvider>
      </MainProvider>
    </>
  );
}

export default App;
