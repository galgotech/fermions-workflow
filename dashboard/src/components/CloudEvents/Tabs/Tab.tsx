import classNames from "classnames";
import { BsX } from "react-icons/bs";

import { CloudEventModel } from "../states/CloudEventModel"

interface Props {
  cloudEvent: CloudEventModel;
  active: boolean;
  open: () => void;
  close: () => void;
}

export const Tab = ({ cloudEvent, active, open, close }: Props) => {
  return (
    <li className="nav-item" style={{"width": "14rem"}}>
      <span
        className={classNames("nav-link link-dark link-underline-opacity-0 link-underline-opacity-0-hover", {
            active: active,
        })}
        role="button"
        onClick={() => open()}
      >
        <div className="d-flex">
          <span className="flex-grow-1">
            {cloudEvent.name}{"  "}
          </span>
          <span onClick={(e) => {
            e.stopPropagation();
            close();
          }}>
            <BsX />
          </span>
        </div>
      </span>
    </li>
  );
}
