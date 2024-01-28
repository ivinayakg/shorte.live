import { useSearchParams } from "react-router-dom";

export default function RedirectPage() {
  let [searchParams] = useSearchParams();

  console.log(searchParams);

  return (
    <div>
      <h1>Redirect Page</h1>
    </div>
  );
}
