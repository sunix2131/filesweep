import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it } from 'vitest';
import { App } from './App';

describe('App', () => {
  it('renders primary navigation', () => {
    render(<MemoryRouter><App /></MemoryRouter>);
    expect(screen.getByText('FileSweep')).toBeInTheDocument();
    expect(screen.getByText('Дубликаты')).toBeInTheDocument();
  });
});
