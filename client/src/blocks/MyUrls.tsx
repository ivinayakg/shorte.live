import { useLocation } from "react-router-dom";
import { useMain } from "@/components/main-provider";
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

export default MyUrls;
