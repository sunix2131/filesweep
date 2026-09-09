import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { useAppStore } from '../stores/appStore';
import type { Settings } from '../types/filesweep';
import { Notice } from '../components/Notice';
import { errorMessage } from '../lib/errors';

const fallback: Settings = { language: 'ru', theme: 'system', includeHiddenFiles: false, followSymlinks: false, maxHashWorkers: 4, maxPreviewFileSizeMb: 20, scanExcludedFolderNames: [], recentFolders: [] };

export function SettingsPage() {
  const { t, setSettings } = useAppStore();
  const [cfg, setCfg] = useState<Settings>(fallback);
  const [error, setError] = useState('');
  const [saved, setSaved] = useState(false);
  const [busy, setBusy] = useState(false);
  useEffect(() => {
    let active = true;
    void api.getSettings().then((next) => { if (active) setCfg(next); }).catch((reason) => { if (active) setError(errorMessage(reason)); });
    return () => { active = false; };
  }, []);
  const save = async () => {
    setBusy(true);
    setError('');
    setSaved(false);
    try {
      const next = await api.updateSettings(cfg);
      setCfg(next);
      setSettings(next);
      setSaved(true);
    } catch (reason) {
      setError(errorMessage(reason));
    } finally {
      setBusy(false);
    }
  };
  return <section><header className="pageHeader"><h1>{t('settings')}</h1><button disabled={busy} onClick={() => void save()}>{t('save')}</button></header>
    {error && <Notice>{error}</Notice>}
    {saved && <Notice kind="success">{t('settingsSaved')}</Notice>}
    <div className="settingsGrid">
      <label>{t('language')}<select value={cfg.language} onChange={(e) => setCfg({ ...cfg, language: e.target.value as Settings['language'] })}><option value="ru">Русский</option><option value="en">English</option></select></label>
      <label>{t('theme')}<select value={cfg.theme} onChange={(e) => setCfg({ ...cfg, theme: e.target.value as Settings['theme'] })}><option value="system">System</option><option value="light">Light</option><option value="dark">Dark</option></select></label>
      <label><input type="checkbox" checked={cfg.includeHiddenFiles} onChange={(e) => setCfg({ ...cfg, includeHiddenFiles: e.target.checked })} />{t('hidden')}</label>
      <label>{t('workers')}<input type="number" value={cfg.maxHashWorkers} onChange={(e) => setCfg({ ...cfg, maxHashWorkers: Number(e.target.value) })} /></label>
      <label>{t('previewSize')}<input type="number" value={cfg.maxPreviewFileSizeMb} onChange={(e) => setCfg({ ...cfg, maxPreviewFileSizeMb: Number(e.target.value) })} /></label>
    </div>
  </section>;
}
