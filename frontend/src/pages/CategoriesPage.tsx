import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { bytes } from '../lib/format';
import { useAppStore } from '../stores/appStore';
import { Empty } from '../components/Empty';
import { Notice } from '../components/Notice';
import { errorMessage } from '../lib/errors';

export function CategoriesPage() {
  const { t, activeScan } = useAppStore();
  const [items, setItems] = useState<Array<{ category: string; filesCount: number; totalSizeBytes: number }>>([]);
  const [error, setError] = useState('');
  useEffect(() => {
    if (!activeScan) {
      setItems([]);
      return;
    }
    let active = true;
    void api.categories(activeScan.id)
      .then((result) => { if (active) { setItems(result); setError(''); } })
      .catch((reason) => { if (active) setError(errorMessage(reason)); });
    return () => { active = false; };
  }, [activeScan]);
  const maxSize = Math.max(...items.map((item) => item.totalSizeBytes), 1);
  return <section><header className="pageHeader"><h1>{t('categories')}</h1></header>
    {error && <Notice>{error}</Notice>}
    {items.length ? <div className="categoryGrid">{items.map((c) => <div className="categoryTile" key={c.category}>
      <span className="categoryIcon" data-category={c.category} />
      <span className="categoryBody">
        <strong>{t(categoryLabelKey(c.category))}</strong>
        <span>{c.filesCount} {t('filesCount')}</span>
        <span>{bytes(c.totalSizeBytes)}</span>
        <i><b style={{ width: `${Math.max(4, (c.totalSizeBytes / maxSize) * 100)}%` }} /></i>
      </span>
    </div>)}</div> : <Empty />}
  </section>;
}

function categoryLabelKey(category: string) {
  const labels: Record<string, 'categoryImages' | 'categoryVideos' | 'categoryAudio' | 'categoryDocuments' | 'categorySpreadsheets' | 'categoryPresentations' | 'categoryArchives' | 'categoryInstallers' | 'categoryCode' | 'categoryFonts' | 'categoryDiskImages' | 'categoryOther'> = {
    images: 'categoryImages', videos: 'categoryVideos', audio: 'categoryAudio', documents: 'categoryDocuments',
    spreadsheets: 'categorySpreadsheets', presentations: 'categoryPresentations', archives: 'categoryArchives',
    installers: 'categoryInstallers', code: 'categoryCode', fonts: 'categoryFonts', disk_images: 'categoryDiskImages', other: 'categoryOther'
  };
  return labels[category] ?? 'categoryOther';
}
