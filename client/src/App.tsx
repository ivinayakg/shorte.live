import { ThemeProvider } from "@/components/theme-provider";
import "./App.css";
import { createBrowserRouter, Outlet, RouterProvider } from "react-router-dom";
import { MainProvider } from "@/components/main-provider";
import { Toaster } from "@/components/ui/toaster";
import { AuthComponent } from "@/utils/auth";
import MyUrls from "@/blocks/MyUrls";
import Header from "@/blocks/Header";
import CreateShort from "@/blocks/CreateShort";
import NotFound from "@/blocks/NotFound";

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
      {
        path: "/not-found/redirect",
        element: <NotFound />,
      },
    ],
  },
  {
    path: "/auth/",
    element: <AuthComponent />,
  },
]);

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
