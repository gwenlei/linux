配置无密码登陆
ssh-keygen
ssh-copy-id -i ~/.ssh/id_rsa.pub root@192.168.122.171
自动填充私钥密码
ssh-agent bash
ssh-add ~/.ssh/id_rsa
手工添加已知主机
ssh-keyscan web3 web4 web5 >> ~/.ssh/known_hosts


https://sysadmincasts.com/episodes/47-zero-downtime-deployments-with-ansible-part-4-4
ssh-keyscan web3 web4 web5 >> .ssh/known_hosts

time ansible-playbook e45-ssh-addkey.yml --ask-pass

ansible all -m ping
ansible web1 -m setup -a "filter=ansible_hostname"
ansible web1 -m setup -a "filter=ansible_eth1"
ansible web1 -m apt -a "name=ntp state=installed"
ab -n 10000 -c 25 http://localhost:8080/
建立互信：减少每次输入密码的麻烦
  ansible all -m copy -a "src=/root/.ssh/id_rsa.pub dest=/root" -k
  ansible all -m shell -a "cat /root/id_rsa.pub >> /root/.ssh/authorized_keys"
  ansible all -m shell -a "rm -f /root/id_rsa.pub"

ansible -i hosts  all -m shell -a "cat /root/id_rsa.pub >> /root/.ssh/authorized_keys"

forks 15

# common
- hosts: all
  sudo: yes
  gather_facts: no
  vars:
    app_version: release-0.01
  serial: 1
  serial: 30%

pre_tasks:
  - name: disable server in haproxy
    shell: echo "disable server {{ inventory_hostname }} "
    delegate_to: "{{ item }}"
    with_items: groups.lb
   -name: re-enable nagios alerts
    nagios: action=enable_alerts host={{ inventory_hostname }} services=webserver
    delegate_to: "{{ item }}"
    with_items: groups.monitoring

post_tasks:

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
  - name: sleep for 5 seconds
    shell: /bin/sleep 5
  - name: write our nginx.conf
    action: template src=templates/nginx.conf.j2 dest=/etc/nginx/nginx.conf
  - name: install nginx
    action: apt name=nginx state=installed
  - name: deploy website content
    git: repo=http://github.com/jweissig...git
         dest=/usr/share/nginx/html
         version={{ app_version }}
  - name: clean existing website content
    #shell: rm -f /usr/share/nginx/html/*
    file: path=/usr/share/nginx/html/ state=absent

  handlers:
  - name: restart nginx
    service: name=nginx state=restarted
  - name: restart haproxy
    service: name=haproxy state=restarted

