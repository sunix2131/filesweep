import { useAppStore } from '../stores/appStore';
export function Empty() {
  const t = useAppStore((s) => s.t);
  return <div className="empty">{t('empty')}</div>;
}
