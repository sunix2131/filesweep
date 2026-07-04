import { Component, type ErrorInfo, type ReactNode } from 'react';

type Props = {
  children: ReactNode;
};

type State = {
  message: string;
};

export class ErrorBoundary extends Component<Props, State> {
  state: State = { message: '' };

  static getDerivedStateFromError(error: Error): State {
    return { message: error.message };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error(error, info.componentStack);
  }

  render() {
    if (!this.state.message) {
      return this.props.children;
    }

    return (
      <div className="fatal">
        <h1>FileSweep</h1>
        <p>Interface failed to load.</p>
        <pre>{this.state.message}</pre>
      </div>
    );
  }
}
