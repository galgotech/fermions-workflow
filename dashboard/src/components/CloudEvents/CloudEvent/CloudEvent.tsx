import { CloudEvent as CloudEventSpec } from "cloudevents";
import { BsFillFileCodeFill } from "react-icons/bs";

import { useDispatch } from "../../../store/types";
import { CloudEventModel } from "../states/CloudEventModel";
import { addEventInStream } from "../states/slice";

import { Editor } from "./Editor";
import { Stream } from "./Stream";



interface Props {
  sendEvent: (cloudEvent: CloudEventSpec<any>) => void;
  cloudEvent: CloudEventModel;
  onChange: (cloudEvent: CloudEventModel) => void;
  openEvent: (raw: string) => void;
}

export const CloudEvent = ({ sendEvent, cloudEvent, onChange, openEvent }: Props) => {
  const dispatch = useDispatch();

  return (
    <div
      className="border-start border-end border-bottom p-2 d-flex flex-column"
      style={{ height: "calc(100% - 44px)" }}
    >
      <div className="d-flex">
        <div className="me-auto">
          <span>
            <BsFillFileCodeFill size={20} />
          </span>
          <span className="fs-4 align-middle">
            {cloudEvent.name}
          </span>
        </div>
        <div>
          <button
            type="button"
            className="btn btn-primary m-1"
            onClick={() => {
              const event = cloudEvent.getCloudEvent();
              dispatch(addEventInStream({name: "send", raw: event.toString()}));
              sendEvent(event);
            }}
          >
            Send
          </button>
        </div>
      </div>
      <div className="d-flex flex-column" style={{ height: "100%" }}>
        <div style={{ height: "40%", minHeight: "200px" }}>
          <Editor
            value={cloudEvent.raw}
            onChange={(value: string) => {
              cloudEvent.raw = value;
              onChange(cloudEvent);
            }}
          />
        </div>
        <hr className="mt-1 mb-1" />
        <div className="flex-grow-1" style={{ height: "60%" }}>
          <Stream openEvent={openEvent} />
        </div>
      </div>
    </div>
  );
};
