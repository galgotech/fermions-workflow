import { Centrifuge } from 'centrifuge';
import { CloudEvent } from "cloudevents";
import { Subscription } from 'rxjs';

import { EventBusSrv } from './event';

export class Workflow {
  private subs = new Subscription();
  private bus = new EventBusSrv();
  private centrifuge?: Centrifuge;

  init (url: string) {
    this.centrifuge = new Centrifuge(url, {
      protocol: 'json',
      debug: true,
    });
    this.centrifuge.on('connecting', function(ctx) {
      console.log(ctx, 'connecting');
      // do whatever you need in case of connecting to a server
    });
    this.centrifuge.on('connected', function(ctx) {
      console.log(ctx, 'connected');
      // now client connected to Centrifugo and authenticated.
    });
    this.centrifuge.on('disconnected', function(ctx) {
      console.log(ctx, 'disconnected');
      // do whatever you need in case of disconnect from server
    });
    // centrifuge.disconnect();
    this.centrifuge.connect();
  }

  connect(on: (event: CloudEvent) => void, source: string, subject?: string) {
    const eventsSub = this.centrifuge?.newSubscription(source);
    eventsSub?.on('publication', (ctx) => {
      const event = new CloudEvent(ctx.data);
      this.bus.publish(event);
    });
    eventsSub?.subscribe();

    this.subs.add(this.bus.subscribe(on, source, subject));
  }

  publish(event: CloudEvent<any>) {
    this.centrifuge?.publish('fermions-worklow-ui', event).then(function(res) {
      console.log('successfully published');
    }, function(err) {
      console.log('publish error', err);
    });

    // this.centrifuge?.rpc("publish", event).then(function(res) {
    //     console.log('rpc result', res);
    // }, function(err) {
    //     console.log('rpc error', err);
    // });
  }
}

const workflow = new Workflow()

export const getWorkflow = () => {
    return workflow;
};
