export function Notice({ children, kind = 'error' }: { children: string; kind?: 'error' | 'success' }) {
  return <div className={`notice notice-${kind}`} role={kind === 'error' ? 'alert' : 'status'}>{children}</div>;
}
