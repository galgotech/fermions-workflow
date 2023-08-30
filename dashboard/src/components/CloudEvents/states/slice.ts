import { PayloadAction, createSlice } from '@reduxjs/toolkit';

import { CloudEventModel } from './CloudEventModel';
import { CloudEventStreamModel } from './CloudEventStreamModel';

function openTab(state: any, action: PayloadAction<CloudEventModel>) {
  const cloudEvent = state.eventsOpen.find((e: CloudEventModel) => e === action.payload);
  if (!cloudEvent) {
    state.eventsOpen.push(action.payload);
  }
}

export interface CoudEventsState {
  current: CloudEventModel;
  events: CloudEventModel[];
  eventsOpen: CloudEventModel[];
  stream: CloudEventStreamModel[];
}

export const cloudEventSlice = createSlice({
  name: 'coudEvents',
  initialState: (): CoudEventsState => {
    const event = new CloudEventModel("New cloud event");
    return {
      current: event,
      events: [event],
      eventsOpen: [event],
      stream: [],
    };
  },
  reducers: {
    open: openTab,
    close: (state, action: PayloadAction<CloudEventModel>) => {
      for (let i = 0; i < state.eventsOpen.length; i++) {
        if (state.eventsOpen[i] === action.payload) {
          if (state.eventsOpen.length -1 !== 0) {
            state.eventsOpen.splice(i, 1);
          }
          break;
        }
      }
    },
    setCurrent: (state, action: PayloadAction<CloudEventModel>) => {
      openTab(state, action);
      state.current = action.payload;
    },
    setEvents: (state, action: PayloadAction<CloudEventModel[]>) => {
      state.events = action.payload;
    },
    newCloudEvent: (state, action: PayloadAction<string | undefined>) => {
      const newEvent = new CloudEventModel("New cloud event", action.payload);
      state.events.push(newEvent);
      state.eventsOpen.push(newEvent);
      state.current = newEvent;
    },
    update: (state, action: PayloadAction<CloudEventModel>) => {
      const cloudEvent = state.events.find((e) => e === action.payload);
      if (!cloudEvent) {
        return;
      }
      cloudEvent.name = action.payload.name;
      cloudEvent.raw = action.payload.raw;
    },
    addEventInStream: (state, action: PayloadAction<{name: string, raw: string}>) => {
      const data = JSON.parse(action.payload.raw);
      delete data["time"]
      delete data["id"]
      const event = new CloudEventModel(action.payload.name, JSON.stringify(data, undefined, 2));
      state.stream.push(new CloudEventStreamModel("send", event));
    },
  },
});

export const {
  open,
  close,
  setCurrent,
  update,
  newCloudEvent,
  addEventInStream,
} = cloudEventSlice.actions;

export default cloudEventSlice.reducer;
