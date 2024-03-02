import { useLocation, useNavigate } from "react-router-dom";
import { useMain } from "@/components/main-provider";
import { useEffect, useRef, useState } from "react";
import fetch from "@/utils/axios";
import { useToast } from "@/components/ui/use-toast";
import UpdateURLModal from "@/components/UpdateURLModal";
import DeleteURLModal from "@/components/DeleteURLModal";
import URLCard from "@/components/URLCard";
import { HeadingOne } from "@/components/typography";

function MyUrls() {
  const { userState } = useMain();
  const [urlsData, setUrlsData] = useState([]);
  const [urlObj, setUrlObj] = useState(null);
  const editRef = useRef<HTMLButtonElement>();
  const deleteRef = useRef<HTMLButtonElement>();
  const { toast } = useToast();
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    if (userState.login && location.pathname === "/my-urls") {
      (async () => {
        const res = await fetch.get("/url/all", {
          withCredentials: true,
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
    }else{
      navigate("/")
    }
  }, [location.pathname, userState.login]);

  return (
    <div className="py-5">
      <UpdateURLModal urlObj={urlObj} ref={editRef} />
      <DeleteURLModal urlObj={urlObj} ref={deleteRef} />
      <HeadingOne className="font-bold mb-3">Your URLs</HeadingOne>
      <div className="flex flex-col justify-center items-center gap-4">
        {urlsData?.map((url: any) => {
          const date = new Date(url.expiry * 1000);
          return (
            <URLCard
              key={url._id}
              shortUrl={url.short}
              destination={url.destination}
              onEditClick={() => {
                setUrlObj(url);
                editRef.current?.click();
              }}
              onDeleteClick={() => {
                setUrlObj(url);
                deleteRef.current?.click();
              }}
              expiry={date}
            />
          );
        })}
      </div>
    </div>
  );
}

export default MyUrls;
