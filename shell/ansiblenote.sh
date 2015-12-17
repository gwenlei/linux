https://sysadmincasts.com/episodes/47-zero-downtime-deployments-with-ansible-part-4-4
ssh-keyscan web3 web4 web5 >> .ssh/known_hosts

ansible-playbook e45-ssh-addkey.yml --ask-pass

ansible all -m ping
ansible web1 -m setup -a "filter=ansible_hostname"
ansible web1 -m setup -a "filter=ansible_eth1"
ab -n 10000 -c 25 http://localhost:8080/

# common
- hosts: all
  sudo: yes
  gather_facts: no

  tasks:
  - name:
    template:
    notify: restart nginx
  - name: enable haproxy
    lineinfile: dest=/etc/default/haproxy regexp="^ENABLED" line="ENABLED=1"
    notify: restart haproxy
  - name: install git
    apt: name=git state=installed update_cache=yes
  - name: install haproxy and socat
    apt: pkg={{ item }} state=latest
    with_items:
    - haproxy
    - socat
  - name: install mysql
    yum: name={{ item }} state=latest
    with_items:
      - mysql-server
      - mysql
      - mysql-deve

  handlers:
  - name: restart nginx
    service: name=nginx state=restarted
  - name: restart haproxy
    service: name=haproxy state=restarted

