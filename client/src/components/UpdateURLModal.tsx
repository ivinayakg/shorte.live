import { useNavigate } from "react-router-dom";
import { UserState } from "@/components/main-provider";
import { forwardRef } from "react";
import fetch from "@/utils/axios";
import { useToast } from "@/components/ui/use-toast";
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

function UpdateURL(
  {
    urlObj,
    userState,
  }: {
    urlObj: any;
    userState: UserState;
  },
  ref: any
) {
  const side = "right";
  const date = new Date(urlObj?.expiry * 1000);
  // const customShort = urlObj?.short.split(
  //   import.meta.env.VITE_REDIRECT_URL_BASE + "/"
  // )[1];
  const navigate = useNavigate();
  const { toast } = useToast();
  const updateURLForm = async (e: any) => {
    e.preventDefault();
    const formData = new FormData(e.target);

    const request = {
      destination: formData.get("destination"),
      short: formData.get("short"),
      expiryStr: formData.get("expiry"),
      expiry: 0,
    };

    if (request.short == "") request.short = null;
    if (request.expiryStr == "") request.expiryStr = null;

    if (request.expiryStr && typeof request.expiryStr == "string") {
      request.expiry = Number(
        (new Date(request.expiryStr).getTime() / 1000).toFixed()
      );
    }

    const res = await fetch.patch(`/url/${urlObj._id}`, request, {
      headers: { Authorization: `Bearer ${userState.token}` },
    });
    if (res.status === 204) {
      toast({
        title: "Short URL updated Successfully",
        duration: 2000,
      });
      setTimeout(() => {
        navigate(0);
      }, 2000);
    }
  };

  const defaultTimeInputValue = urlObj ? date?.toISOString().slice(0, 16) : "";

  return (
    <>
      <Sheet key={side}>
        <SheetTrigger asChild>
          <Button ref={ref} style={{ display: "none" }}></Button>
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
                defaultValue={urlObj?.destination}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-rows-2 items-center gap-4">
              <Label className="text-left">
                Custom Short - <i>{urlObj?.short}</i>
              </Label>
              <Input
                name="short"
                className="col-span-3"
              />
            </div>
            <div className="grid grid-rows-2 items-center gap-4">
              <Label className="text-left">
                Expiry - <i>{date.toLocaleString()}</i>
              </Label>
              <Input
                name="expiry"
                className="col-span-3"
                type="datetime-local"
                defaultValue={defaultTimeInputValue}
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

const UpdateURLModal = forwardRef(UpdateURL);

export default UpdateURLModal;
