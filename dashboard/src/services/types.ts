import { Unsubscribable, Observable } from 'rxjs';
import { CloudEvent } from "cloudevents";

export interface BusEventHandler {
    (event: CloudEvent): void;
}

export interface EventFilterOptions {
    onlyLocal: boolean;
}

export interface EventFilterOptions {
    onlyLocal: boolean;
}  

export interface EventBus {
    /**
     * Publish single event
     */
    publish(event: CloudEvent): void;
  
    /**
     * Get observable of events
     */
    getStream(eventType: string): Observable<CloudEvent>;
  
    /**
     * Subscribe to an event stream
     *
     * This function is a wrapper around the `getStream(...)` function
     */
    subscribe(handler: BusEventHandler, source: string, subject?: string): Unsubscribable;
  
    /**
     * Remove all event subscriptions
     */
    removeAllListeners(): void;
  
    /**
     * Returns a new bus scoped that knows where it exists in a heiarchy
     *
     * @internal -- This is included for internal use only should not be used directly
     */
    newScopedBus(key: string, filter: EventFilterOptions): EventBus;
  }
  