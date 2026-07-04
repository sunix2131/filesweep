import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { bytes } from '../lib/format';
import { useAppStore } from '../stores/appStore';
import { Empty } from '../components/Empty';

export function CategoriesPage() {
  const { t, activeScan } = useAppStore();
  const [items, setItems] = useState<Array<{ category: string; filesCount: number; totalSizeBytes: number }>>([]);
  useEffect(() => { if (activeScan) api.categories(activeScan.id).then(setItems); }, [activeScan]);
  const maxSize = Math.max(...items.map((item) => item.totalSizeBytes), 1);
  return <section><header className="pageHeader"><h1>{t('categories')}</h1></header>
    {items.length ? <div className="categoryGrid">{items.map((c) => <button className="categoryTile" key={c.category}>
      <span className="categoryIcon" data-category={c.category} />
      <span className="categoryBody">
        <strong>{categoryLabel(c.category)}</strong>
        <span>{c.filesCount} файлов</span>
        <span>{bytes(c.totalSizeBytes)}</span>
        <i><b style={{ width: `${Math.max(4, (c.totalSizeBytes / maxSize) * 100)}%` }} /></i>
      </span>
    </button>)}</div> : <Empty />}
  </section>;
}

function categoryLabel(category: string) {
  const labels: Record<string, string> = {
    images: 'Изображения',
    videos: 'Видео',
    audio: 'Аудио',
    documents: 'Документы',
    spreadsheets: 'Таблицы',
    presentations: 'Презентации',
    archives: 'Архивы',
    installers: 'Установщики',
    code: 'Код',
    fonts: 'Шрифты',
    disk_images: 'Диск-образы',
    other: 'Другое'
  };
  return labels[category] ?? category;
}
