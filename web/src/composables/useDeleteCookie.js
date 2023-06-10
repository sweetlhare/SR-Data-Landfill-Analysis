import { useSetCookie } from './useSetCookie';

export function useDeleteCookie(name) {
  useSetCookie(name, '', {
    'max-age': -1,
  });
}
