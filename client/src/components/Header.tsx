import { Button } from "@/components/ui/button";
import { ModeToggle } from "@/components/mode-toggle";
import { useNavigate, Link } from "react-router-dom";
import { useMain } from "@/components/main-provider";
import { setInLocalStorage } from "@/utils/localstorage";
import { AvatarImage, Avatar, AvatarFallback } from "@/components/ui/avatar";

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
    <div className="min-w-full flex justify-between items-center py-5 flex-col sm:flex-row gap-4">
      <Link to="/">
        <h1 className="text-3xl">shorte.live</h1>
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

export default Header;
