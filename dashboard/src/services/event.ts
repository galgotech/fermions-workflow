import { CloudEvent } from "cloudevents";
import EventEmitter from 'eventemitter3';
import { Observable, Unsubscribable, Subscriber } from 'rxjs';

import { BusEventHandler, EventBus, EventFilterOptions } from './types';

export class EventBusSrv implements EventBus {
  private emitter: EventEmitter;
  private subscribers = new Map<Function, Subscriber<CloudEvent>>();

  constructor() {
    this.emitter = new EventEmitter();
  }

  publish(event: CloudEvent): void {
    this.emitter.emit(event.source + (event.subject || ''), event);
  }

  subscribe(handler: BusEventHandler, source: string, subject?: string): Unsubscribable {
    return this.getStream(source + (subject || '')).subscribe({ next: handler });
  }

  getStream(eventType: string): Observable<CloudEvent> {
    return new Observable((observer) => {
      const handler = (event: CloudEvent) => {
        observer.next(event);
      };

      this.emitter.on(eventType, handler);
      this.subscribers.set(handler, observer);

      return () => {
        this.emitter.off(eventType, handler);
        this.subscribers.delete(handler);
      };
    });
  }

  newScopedBus(key: string, filter?: EventFilterOptions): EventBus {
    throw new Error("not implemented");
  // return new ScopedEventBus([key], this, filter);
  }

  removeAllListeners() {
    this.emitter.removeAllListeners();
    for (const [key, sub] of this.subscribers) {
      sub.complete();
      this.subscribers.delete(key);
    }
  }
}
