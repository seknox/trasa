module.exports = {
  docs: [
    // 'overview',
    // 'concepts',
    // 'glossary',
    {
      'Getting Started': [
        'getting-started/overview',
        'getting-started/concepts',
        'getting-started/glossary',
        'getting-started/signup-or-install',

        'getting-started/how-to',
      ],
      // Tutorial: ['tutorial/tutorial'],
      Install: ['install/installation', 'install/initial-setup'],
      QuickStart: ['quickstarts/quickstart-cloud-ssh'],
      Providers: [{ users: ['providers/users/ldap/ldap'] }],
      Users: [
        'users/users',
        'users/creating-updating-users',
        // 'users/trasaIDP/trasa-idp',
        // 'users/ldap/ldap',
        // 'users/ad/active-directory',
        // 'users/freeIPA/free-ipa',
      ],
      Services: [
        'services/introduction',
        'services/privilege',
        'services/http/http-service',
        'services/rdp/rdp-service',
        'services/ssh/ssh-service',
        'services/radius/radius-server',
      ],
      Policies: ['policies/policy-introduction', 'policies/basic-policy', 'policies/device-policy'],
      'Access Map': ['access-map/access-mapping'],
      'Access Proxy': ['access-proxy/introduction', 'access-proxy/sshproxy'],
      'Two Factor Authentication': [
        'native-tfa/two-factor-authentication',
        'native-tfa/windows/windows-two-factor-authentication',
        'native-tfa/linux-two-factor-authentication',
      ],
      'Secret Vault': ['providers/secret-vault/index'],

      'System Configurations': ['system/index', 'system/config-reference'],
    },
  ],
  guides: [
    'guides/getting-started',
    {
      Account: ['guides/user/account/password-setup'],
      Device: [
        'guides/user/device/enrol-2fa-device',
        'guides/user/device/install-trasa-device-agent',
      ],
      Access: [
        'guides/user/access/dashboard-login',
        'guides/user/access/adhoc-access',
        'guides/user/access/tfa',
        'guides/user/access/ssh-connection-via-proxy',
        'guides/user/access/rdp-connection-via-proxy',
        'guides/user/access/https-connection-via-proxy',
      ],
    },
  ],
};
