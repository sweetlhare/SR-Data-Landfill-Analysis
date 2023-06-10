import { useGetCookie } from './useGetCookie';

export function useGetSessionId() {
  return useGetCookie('reoLandfillsSessionId');
}
