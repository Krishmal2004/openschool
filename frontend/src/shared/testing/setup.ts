class ResizeObserverMock {
  observe() {}
  unobserve() {}
  disconnect() {}
}

if (typeof window.ResizeObserver === "undefined") {
  window.ResizeObserver = ResizeObserverMock as unknown as typeof ResizeObserver;
}
