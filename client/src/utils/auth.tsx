import { useNavigate, useSearchParams } from "react-router-dom";
import { setInLocalStorage } from "./localstorage";
import { useEffect } from "react";

function AuthComponent() {
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

export { AuthComponent };
