import { BsArrowDownShort, BsArrowUpShort } from "react-icons/bs";

import { useSelector } from '../../../store/types';

interface Props {
  openEvent: (raw: string) => void;
}

export const Stream = ({ openEvent }: Props) => {
  const cloudEvents = useSelector((state) => state.cloudEvents);

  return (
    <>
      <div className="fs-5">Cloud events</div>
        {cloudEvents.stream.map((stream, i: number) => (
          <div
            key={i}
            className="fw-light"
            onClick={() => openEvent(stream.cloudEvent.raw)}
          >
            <BsArrowDownShort
              style={{visibility: stream.direction === "send" ? "hidden" : "visible"}}
            />
          <BsArrowUpShort
            style={{visibility: stream.direction === "send" ? "visible" : "hidden"}}
          />
          {" "}
          {stream.cloudEvent.raw}
        </div>
      ))}
    </>
  );
}
