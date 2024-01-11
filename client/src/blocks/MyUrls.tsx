import { useLocation, useNavigate } from "react-router-dom";
import { useMain } from "@/components/main-provider";
import { useEffect, useRef, useState } from "react";
import fetch from "@/utils/axios";
import { useToast } from "@/components/ui/use-toast";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import UpdateURLModal from "@/blocks/UpdateURLModal";

function MyUrls() {
  const { userState } = useMain();
  const [urlsData, setUrlsData] = useState([]);
  const [urlObj, setUrlObj] = useState(null);
  const editRef = useRef<HTMLButtonElement>();
  const deleteRef = useRef<HTMLButtonElement>();
  const { toast } = useToast();
  const Headings = [
    "Serial",
    "Short",
    "Destination",
    "Expiry",
    "ID",
    "Edit",
    "Delete",
  ];
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
      <UpdateURLModal urlObj={urlObj} userState={userState} ref={editRef} />
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
                    <Button
                      onClick={() => {
                        setUrlObj(url);
                        editRef.current?.click();
                      }}
                      variant="outline"
                      className="font-bold"
                    >
                      Edit
                    </Button>
                  </TableCell>
                  <TableCell className="font-medium text-left">
                    <Button
                      onClick={() => {
                        setUrlObj(url);
                        deleteRef.current?.click();
                      }}
                      variant="destructive"
                      className="font-bold"
                    >
                      Delete
                    </Button>
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
