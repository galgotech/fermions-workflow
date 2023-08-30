import { CloudEvent as CloudEventSpec } from "cloudevents";
import { useEffect, useState } from "react";
import { BsPlus } from "react-icons/bs";

import { getWorkflow } from "../../services/workflow";
import { useDispatch, useSelector } from '../../store/types';

import { CloudEvent } from './CloudEvent/CloudEvent';
import { EventName } from './Nav/EventName';
import { Tabs } from './Tabs/Tabs';
import { CloudEventModel } from './states/CloudEventModel';
import { addEventInStream, close, newCloudEvent, setCurrent, update } from './states/slice';


export const CloudEvents = () => {
  const dispatch = useDispatch();
  const { events, eventsOpen, current } = useSelector((state) => state.cloudEvents);

  const [render, setRender] = useState(0);

  const updateCloudEvent = (cloudEvent: CloudEventModel) => {
    dispatch(update(cloudEvent));
    setRender(render + 1);
  };

  const closeTab = (cloudEvent: CloudEventModel) => {
    dispatch(close(cloudEvent));
    setRender(render + 1);
  };

  const openTab = (cloudEvent: CloudEventModel) => {
    dispatch(setCurrent(cloudEvent));
    setRender(render + 1);
  };

  const sendEvent = (cloudEvent: CloudEventSpec) => {
    getWorkflow().publish(cloudEvent);
  };

  useEffect(() => {
    // getWorkflow().connect((event: CloudEvent2) => {
    //   console.log('-- connect', event);
    // }, "event-source");
    getWorkflow().connect((event: CloudEventSpec<any>) => {
      dispatch(addEventInStream({ name: "receive", raw: event.toString() }));
    }, "event-source2");

  }, [getWorkflow]);

  return (
    <div className="d-flex p-1" style={{ height: "100%" }}>
      <div style={{ width: "16rem" }}>
        <div>
          <BsPlus />
          <span
            className="fw-light link-dark link-underline-opacity-0 link-underline-opacity-0-hover align-middle"
            role="button"
            onClick={() => dispatch(newCloudEvent())}
          >
            new cloud event
          </span>
        </div>
        <hr className="mt-1 mb-1" />
        <div className="ms-2 mt-2">
          {events.map((cloudEvent: CloudEventModel, i: number) => (
            <EventName
              key={i}
              cloudEvent={cloudEvent}
              open={(cloudEvent: CloudEventModel) => openTab(cloudEvent)}
              save={(cloudEvent: CloudEventModel) => updateCloudEvent(cloudEvent)}
            />
          ))}
        </div>
      </div>

      <div className="ms-4 flex-grow-1">
        <Tabs
          eventsOpen={eventsOpen}
          open={(cloudEvent) => openTab(cloudEvent)}
          close={(cloudEvent) => closeTab(cloudEvent)}
        />
        <CloudEvent
          cloudEvent={current}
          onChange={(cloudEvent) => updateCloudEvent(cloudEvent)}
          sendEvent={sendEvent}
          openEvent={(raw) => dispatch(newCloudEvent(raw))}
        />
      </div>
    </div>
  );
};
