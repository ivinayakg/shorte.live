import { useLocation, useNavigate } from "react-router-dom";
import { UserState, useMain } from "@/components/main-provider";
import { useEffect, useState } from "react";
import fetch from "@/utils/axios";
import { useToast } from "@/components/ui/use-toast";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";

function UpdateURL({
  urlObj,
  userState,
}: {
  urlObj: any;
  userState: UserState;
}) {
  const side = "right";
  const date = new Date(urlObj.expiry);
  const customShort = urlObj.short.split(
    import.meta.env.VITE_BASE_URL + "/"
  )[1];
  const navigate = useNavigate();

  const updateURLForm = async (e: any) => {
    e.preventDefault();
    const formData = new FormData(e.target);

    const request = {
      destination: formData.get("destination"),
      short: formData.get("short"),
      expiry: formData.get("expiry"),
    };

    if (request.short == "") request.short = null;
    if (request.expiry == "") request.expiry = null;

    if (request.expiry && typeof request.expiry == "string")
      request.expiry = new Date(request.expiry).toISOString();

    const res = await fetch.patch(`/url/${urlObj._id}`, request, {
      headers: { Authorization: `Bearer ${userState.token}` },
    });
    if (res.status === 204) {
      navigate(0);
    }
  };

  return (
    <>
      <Sheet key={side}>
        <SheetTrigger asChild>
          <Button variant="outline" className="font-bold">
            Edit
          </Button>
        </SheetTrigger>
        <SheetContent side={side}>
          <SheetHeader>
            <SheetTitle>Edit URL</SheetTitle>
            <SheetDescription>
              Make changes to your URL here. Click save when you're done.
            </SheetDescription>
          </SheetHeader>
          <form className="grid gap-4 py-4" onSubmit={updateURLForm}>
            <div className="grid grid-rows-2 items-center gap-4">
              <Label className="text-left">Destination</Label>
              <Input
                name="destination"
                defaultValue={urlObj.destination}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-rows-2 items-center gap-4">
              <Label className="text-left">
                Custom Short - <i>{customShort}</i>
              </Label>
              <Input name="short" className="col-span-3" />
            </div>
            <div className="grid grid-rows-2 items-center gap-4">
              <Label className="text-left">
                Expiry - <i>{date.toLocaleString()}</i>
              </Label>
              <Input
                name="expiry"
                className="col-span-3"
                type="datetime-local"
              />
            </div>
            <SheetFooter>
              <SheetClose asChild>
                <Button type="submit">Save changes</Button>
              </SheetClose>
            </SheetFooter>
          </form>
        </SheetContent>
      </Sheet>
    </>
  );
}

function MyUrls() {
  const { userState } = useMain();
  const [urlsData, setUrlsData] = useState([]);
  const { toast } = useToast();
  const Headings = ["Serial", "Short", "Destination", "Expiry", "ID", "Edit"];
  const location = useLocation();
  const navigate = useNavigate();

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

    if (!userState.login) navigate("/");
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
                  <TableCell className="font-medium text-left">
                    <UpdateURL urlObj={url} userState={userState} />
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

export default MyUrls;
