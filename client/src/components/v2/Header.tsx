import { Button } from "@/components/ui/button";
import { ModeToggle } from "@/components/mode-toggle";
import { Link, useLocation } from "react-router-dom";
import { useMain } from "@/components/main-provider";
import { AvatarImage, Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { baseURL } from "@/utils/axios";
import { GitHubLogoIcon } from "@radix-ui/react-icons";

function MobileHeader({ links }: { links: JSX.Element[] }) {
  return (
    <div className="user-navbar-mobile sm:hidden flex justify-around items-center w-11/12 gap-4 bg-secondary p-2 rounded-md fixed bottom-6">
      {links}
    </div>
  );
}

function UserAvatar({ userPicture }: { userPicture: string }) {
  const logoutURL = baseURL + "/user/logout";

  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Avatar className="">
          <AvatarImage src={userPicture} alt="user_profile" />
          <AvatarFallback>U</AvatarFallback>
        </Avatar>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuLabel className="max-w-fit">
          <Link to={logoutURL} className="bg-transparent text-primary">
            Logout
          </Link>
        </DropdownMenuLabel>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function Header() {
  const loginWithGoogleUrl = import.meta.env.VITE_GOOGLE_SIGN_IN;
  const githubRepoUrl = import.meta.env.VITE_GITHUB_REPO;
  const myUrl = "/my-urls";
  const homeUrl = "/";
  const { userState } = useMain();
  const location = useLocation();

  const homeButtonClass =
    location.pathname === "/"
      ? ""
      : "bg-transparent text-primary hover:bg-primary hover:text-secondary";
  const urlsButtonClass =
    location.pathname === "/my-urls"
      ? ""
      : "bg-transparent text-primary hover:bg-primary hover:text-secondary";

  const links = [
    <Link to={homeUrl}>
      <Button variant="default" className={homeButtonClass}>
        Home
      </Button>
    </Link>,
    <Link to={myUrl}>
      <Button variant="default" className={urlsButtonClass}>
        My URLs
      </Button>
    </Link>,
    <UserAvatar userPicture={userState.picture} />,
  ];

  return (
    <>
      <div className="min-w-full flex justify-between items-center p-4 gap-4 border-b-2">
        <Link to="/" className="flex items-center gap-4">
          <span className="w-7 aspect-square rounded-full block bg-primary"></span>
          <h1 className="text-3xl mb-[5px]">shorte.live</h1>
        </Link>
        <div className="flex justify-center items-center gap-4">
          <Link to={githubRepoUrl}>
            <Button className="gap-2">
              Star on Github
              <GitHubLogoIcon className="w-5 h-5" />
            </Button>
          </Link>
          {userState.login ? (
            <div className="user-navbar-tablet hidden sm:flex justify-center items-center gap-4">
              {links}
            </div>
          ) : (
            <Link to={loginWithGoogleUrl}>
              <Button className="gap-2">
                Login <img src="/google.svg" alt="" className="w-5 h-5" />
              </Button>
            </Link>
          )}
          <ModeToggle />
        </div>
      </div>
      {userState.login && <MobileHeader links={links} />}
    </>
  );
}

export default Header;
