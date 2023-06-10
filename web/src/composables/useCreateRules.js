import { required, helpers } from '@vuelidate/validators';

export function useCreateRules() {
  const requiredWithMessage = helpers.withMessage('Обязательное поле', required);
  const rules = {
    name: { requiredWithMessage },
    city: { requiredWithMessage },
    address: { requiredWithMessage },
    latitude: { requiredWithMessage },
    longitude: { requiredWithMessage },
  };

  return rules;
}
