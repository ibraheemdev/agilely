const authenticityToken = () => {
  const token = document.querySelector('meta[name="csrf-token"]');
  if (token && token instanceof window.HTMLMetaElement) {
    return token.content;
  }
  return null;
}

export default authenticityToken