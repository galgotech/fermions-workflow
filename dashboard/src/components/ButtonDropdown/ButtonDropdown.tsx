import { useEffect, useRef, useState } from "react";
import { BsThreeDots } from "react-icons/bs";


enum Action {
  RENAME,
  DELETE,
}

interface Props {
  rename: Function;
  del: Function;
}

export const ButtonDropdown = ({ rename, del }: Props) => {
  const ref = useRef<HTMLUListElement>(null);
  const [dropdown, setDropdown] = useState(false);

  const click = (action: Action) => {
    if (Action.RENAME === action) {
      rename();
    } else if (Action.DELETE === action) {
      del();
    }
    setDropdown(false);
  };

  useEffect(() => {
    const handler = (event: any) => {
      if (dropdown && ref.current && !ref.current.contains(event.target)) {
        setDropdown(false);
      }
    };
    document.addEventListener("mousedown", handler);
    document.addEventListener("touchstart", handler);
    return () => {
      // Cleanup the event listener
      document.removeEventListener("mousedown", handler);
      document.removeEventListener("touchstart", handler);
    };
   }, [dropdown]);

  return (<>
    <BsThreeDots role="button" onClick={() => setDropdown(true)} />
    {dropdown && (
      <ul ref={ref} className="dropdown-menu" style={{display: "block"}}>
        <li>
          <span
            className="dropdown-item"
            role="button"
            onClick={() => click(Action.RENAME)}
          >
            Rename
          </span>
        </li>
        <li>
          <span
            className="dropdown-item"
            role="button"
            onClick={() => click(Action.DELETE)}
          >
            Delete
          </span>
        </li>
      </ul>
    )}
  </>);
}
