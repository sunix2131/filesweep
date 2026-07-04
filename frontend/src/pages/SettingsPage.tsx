import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { useAppStore } from '../stores/appStore';
import type { Settings } from '../types/filesweep';

const fallback: Settings = { language: 'ru', theme: 'system', includeHiddenFiles: false, followSymlinks: false, maxHashWorkers: 4, maxPreviewFileSizeMb: 20, scanExcludedFolderNames: [], recentFolders: [] };

export function SettingsPage() {
  const { t, setSettings } = useAppStore();
  const [cfg, setCfg] = useState<Settings>(fallback);
  useEffect(() => { api.getSettings().then(setCfg).catch(() => undefined); }, []);
  const save = async () => { const next = await api.updateSettings(cfg); setCfg(next); setSettings(next); };
  return <section><header className="pageHeader"><h1>{t('settings')}</h1><button onClick={save}>{t('save')}</button></header>
    <div className="settingsGrid">
      <label>{t('language')}<select value={cfg.language} onChange={(e) => setCfg({ ...cfg, language: e.target.value as Settings['language'] })}><option value="ru">Русский</option><option value="en">English</option></select></label>
      <label>{t('theme')}<select value={cfg.theme} onChange={(e) => setCfg({ ...cfg, theme: e.target.value as Settings['theme'] })}><option value="system">System</option><option value="light">Light</option><option value="dark">Dark</option></select></label>
      <label><input type="checkbox" checked={cfg.includeHiddenFiles} onChange={(e) => setCfg({ ...cfg, includeHiddenFiles: e.target.checked })} />{t('hidden')}</label>
      <label><input type="checkbox" checked={cfg.followSymlinks} onChange={(e) => setCfg({ ...cfg, followSymlinks: e.target.checked })} />{t('symlinks')}</label>
      <label>{t('workers')}<input type="number" value={cfg.maxHashWorkers} onChange={(e) => setCfg({ ...cfg, maxHashWorkers: Number(e.target.value) })} /></label>
      <label>{t('previewSize')}<input type="number" value={cfg.maxPreviewFileSizeMb} onChange={(e) => setCfg({ ...cfg, maxPreviewFileSizeMb: Number(e.target.value) })} /></label>
    </div>
  </section>;
}
