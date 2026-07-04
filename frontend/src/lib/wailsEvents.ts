type Runtime = {
  EventsOn?: (eventName: string, callback: (...data: unknown[]) => void) => () => void;
};

const runtime = () => (window as unknown as { runtime?: Runtime }).runtime;

export function onWailsEvent<T>(eventName: string, callback: (payload: T) => void): () => void {
  const unsubscribe = runtime()?.EventsOn?.(eventName, (...data: unknown[]) => {
    callback(data[0] as T);
  });
  return unsubscribe ?? (() => undefined);
}
