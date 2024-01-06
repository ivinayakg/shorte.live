import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { ThemeProvider } from "@/components/theme-provider";
import "./App.css";
import { ModeToggle } from "./components/mode-toggle";
import { Input } from "./components/ui/input";
import fetch from "./utils/axios";
import { getFromLocalStorage, setInLocalStorage } from "./utils/localstorage";

function Header() {
  const loginWithGoogle = () => {
    const loginWithGoogleUrl = "";
    history.pushState(null, "", loginWithGoogleUrl);
  };

  return (
    <div className="min-w-full flex justify-between items-center py-5">
      <h1 className="text-3xl">URL Shortner</h1>
      <div className="flex justify-center items-center gap-4">
        <Button className="gap-2" onClick={loginWithGoogle}>
          Login With Google <img src="/google.svg" alt="" />
        </Button>
        <ModeToggle />
      </div>
    </div>
  );
}

function CreateShort() {
  return (
    <div className="w-full flex flex-col justify-center items-center gap-10 createshort">
      <h2 className="text-7xl font-bold">Generate shorten link</h2>
      <div className="flex w-full max-w-5xl items-center space-x-2">
        <Input className="" type="text" placeholder="URL" />
        <Input className="w-2/5" type="text" placeholder="shorten" value={""} />
        <Button className="" type="submit">
          Generate
        </Button>
      </div>
    </div>
  );
}

function Main() {
  // const [userState, setUserState] = useState({});

  useEffect(() => {
    (async () => {
      const token = getFromLocalStorage("userToken");
      if (!token) return;
      const res = await fetch.get("/user/self");
      if (res.status !== 200) {
        setInLocalStorage("userToken", null);
        return;
      }
      const data = res.data;
      console.log(data);
    })();
  }, []);

  return (
    <>
      <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <Header />
        <CreateShort />
      </ThemeProvider>
    </>
  );
}

function App() {
  return <Main />;
}

export default App;
