import { useSelector } from "../../../store/types";
import { CloudEventModel } from "../states/CloudEventModel";

import { Tab } from "./Tab";

interface Props {
  eventsOpen: CloudEventModel[];
  open: (cloudEvent: CloudEventModel) => void;
  close: (cloudEvent: CloudEventModel) => void;
}

export const Tabs = ({ eventsOpen, open, close }: Props) => {
  const cloudEvents = useSelector((state) => state.cloudEvents);

  return (
    <ul className="nav nav-tabs">
      {eventsOpen.map((cloudEvent: CloudEventModel, i: number) => (
        <Tab
          key={i}
          cloudEvent={cloudEvent}
          active={cloudEvent === cloudEvents.current}
          open={() => open(cloudEvent)}
          close={() => close(cloudEvent)}
        />
      ))}
    </ul>
  );
}
