import { CloudEventModel } from "./CloudEventModel";


export class CloudEventStreamModel {
  constructor(
    public direction: string,
    public cloudEvent: CloudEventModel,
  ) {
  }
}