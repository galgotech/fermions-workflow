import { CloudEvent } from "cloudevents";

type typesEvent = object | string | number | boolean | undefined;

export class CloudEventModel {
  public raw: string;

  constructor(
    public name: string,
    raw?: string,
  ) {
    if (raw) {
      this.raw = raw;
    } else {
      this.raw = `{
  "specversion" : "1.0",
  "type" : "com.example.someevent",
  "source" : "/mycontext",
  "datacontenttype" : "application/json",
  "data" : {}
}`;
    }
  }

  setRaw(raw: string) {
    this.raw = JSON.stringify(JSON.parse(raw), null, 2);
  }

  getCloudEvent(): CloudEvent<typesEvent> {
    const event = JSON.parse(this.raw);
    // event["time"] = new Date().toString();
    return new CloudEvent<typesEvent>(event);
  }

  valid(): boolean {
    try {
      this.getCloudEvent();
      return true;
    } catch (e) {
      console.log('invalid event', e);
    }
    return false;
  }

  json(): string {
    return this.raw;
  }
}
