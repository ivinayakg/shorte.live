import { useNavigate } from "react-router-dom";
import { UserState } from "@/components/main-provider";
import { forwardRef } from "react";
import fetch from "@/utils/axios";
import { useToast } from "@/components/ui/use-toast";
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
import { Button } from "@/components/ui/button";

function DeleteURL(
  {
    urlObj,
    userState,
  }: {
    urlObj: any;
    userState: UserState;
  },
  ref: any
) {
  const navigate = useNavigate();
  const { toast } = useToast();
  const deleteURLForm = async () => {
    const res = await fetch.delete(`/url/${urlObj._id}`, {
      headers: { Authorization: `Bearer ${userState.token}` },
    });
    if (res.status === 204) {
      toast({
        title: "Short URL deleted Successfully",
        duration: 2000,
      });
      setTimeout(() => {
        navigate("/my-urls");
      }, 2000);
    }
  };

  return (
    <>
      <AlertDialog>
        <AlertDialogTrigger asChild>
          <Button ref={ref} style={{ display: "none" }}></Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete your
              URL.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={deleteURLForm}>
              Continue
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

const DeleteURLModal = forwardRef(DeleteURL);

export default DeleteURLModal;
