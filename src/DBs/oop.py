# coding: utf-8
import os
from sqlite3 import connect
import sys

db = sys.argv[3].upper() + '.db'

class FileSystem:
    # db = sys.argv[3].upper() + '.db'
    con = connect(db)

    def get_elems(self, src_path):
        path_elems = os.walk(src_path).__next__()
        self.l_dirs = [Ffolder(name = x) for x in path_elems[1]]
        for dir in self.l_dirs:
            dir.src_path = src_path
            dir.get_elems(dir.src_path + os.sep + dir.name)
        self.l_files = [Ffile(name = x) for x in path_elems[2]]
        for f in self.l_files:
            f.src_path = src_path + os.sep + f.name
            f.load()

class Src(FileSystem):

    __slots__ = ('id', 'id_prj', 'id_file')

    def __init__(self, id = None, id_prj = None, id_file = None, num = None, line = None, comment = None, tags = None, todo = False):
        self.id = id
        self.id_prj = id_prj
        self.id_file = id_file
        self.num = num
        if line:
            self.line = line.replace("'", "👆").replace('"', '👇').replace("🧨", "").rstrip()
            if not comment:
                l_comments = self.line.split("📜")
                if len(l_comments) > 1:
                    self.comment = l_comments[-1]
                else:
                    self.comment = None
            else:
                self.comment = comment

            if not todo:
                if line.find("🧨") > -1:
                    self.todo = True
                else:
                    self.todo = False
            else:
                self.todo = todo
        else:
            self.line = line


        self.tags = tags
        

    def save(self):
        
        cur = self.con.cursor()
        cur.execute("""
                    insert into src(id_prj
                                    , id_file
                                    , num
                                    , line
                                    , comment
                                    , tags
                                    , todo)
                            values({id_prj}
                                    , {id_file}
                                    , {num}
                                    , {line}
                                    , {comment}
                                    , {tags}
                                    , {todo})
                    """.format(
                        id_prj = '{id_prj}'.format(id_prj=self.id_prj) if self.id_prj else 'null'
                        , id_file = '{id_file}'.format(id_file=self.id_file) if self.id_file else 'null'
                        , num = '{num}'.format(num=self.num) if self.num else 'null'
                        , line = "'{line}'".format(line=self.line) if self.line else 'null'
                        , comment = "'{comment}'".format(comment=self.comment) if self.comment else 'null'
                        , tags = "'{tags}'".format(tags=self.tags) if self.tags else 'null'
                        , todo = '{todo}'.format(todo=self.todo) if self.todo else 'False'))
        
        cur.execute(f"""
                      update src
                         set line = replace(line, '👆', char(39))
                       where id_file = {self.id_file}
                    """)
        
        cur.execute(f"""
                      update src
                         set line = replace(line, '👇', char(34))
                       where id_file = {self.id_file}
                    """)
        
        self.con.commit()

class Ffile(FileSystem):


    __slots__ = ('id', 'name', 'id_prj', 'id_folder', 'comment', 'tags')

    src_path = None

    l_lines = list()

    def __init__(self, id = None, id_prj = None, id_folder = None, name = None, comment = None, tags = None):
        self.id = id
        self.name = name
        self.id_prj = id_prj
        self.id_folder = id_folder
        self.comment = comment
        self.tags = tags

        if self.id:
            cur = self.con.cursor()
            cur.execute(f"""select id
                                   , id_prj 
                                   , id_file
                                   , num
                                   , line
                                   , comment
                                   , tags
                              from src
                             where id_file = {self.id}
                            """)
            self.l_lines = [Src(*x) for x in cur.fetchall()]


    def download(self, down_path: str, todo: bool = False):
        with open(down_path + os.sep + self.name, encoding='utf-8', mode='w') as f:
            for src in self.l_lines:
                line = src
                if line:
                    if todo:
                        if not todo:
                            f.write(line.replace("👆", "'").replace('👇', '"') + '\n')
                        else:
                            f.write('...' + '\t\t' + self.comment + '\n')
                    else:
                        f.write(line.replace("👆", "'").replace('👇', '"') + '\n')
                else:
                    f.write("\n")                        


    def load(self):
        with open(self.src_path, encoding="utf-8", mode="r") as f: l_lines = f.readlines()
        self.l_lines = [Src(num=i+1, line=l_lines[i], id_file=self.id) for i in range(len(l_lines))]


    def save(self):
        cur = self.con.cursor()
        cur.execute("""
                            insert into obj(id_prj
                                                , id_parent
                                                , name
                                                , comment
                                                , object_type
                                                , tags)
                                        values ({id_prj}
                                                , {id_parent}
                                                , {name}
                                                , {comment}
                                                , 1
                                                , {tags})
                                returning id
                            """.format(id_prj = '{id_prj}'.format(id_prj=self.id_prj) if self.id_prj else 'null'
                                    , id_parent = '{id_parent}'.format(id_parent=self.id_folder) if self.id_folder else 'null'
                                    , name = "'{name}'".format(name=self.name) if self.name else 'null'
                                    , comment = "'{comment}'".format(comment=self.comment) if self.comment else 'null'
                                    , tags = "'{tags}'".format(tags=self.tags) if self.tags else 'null'))

        self.id = cur.fetchone()[0]

        for line in self.l_lines: 
            line.id_prj = self.id_prj
            line.id_file = self.id
            line.save()
        
        self.con.commit()
        

class Ffolder(FileSystem):

    __slots__ = ('id', 'id_prj', 'id_parent', 'name', 'comment', 'tags')

    src_path = None
    
    l_dirs = list()
    l_files = list()
    
    def __init__(self, id = None, id_prj = None, id_parent = None, name = None, comment = None, tags = None):
        self.id = id
        self.id_prj = id_prj
        self.id_parent = id_parent
        self.name = name
        self.comment = comment
        self.tags = tags
        # self.path = os.path.basename(src_path)
        
        if self.id:
            cur = self.con.cursor()
            cur.execute(f"""
                            select id
                                , id_prj
                                , id_parent
                                , name
                                , comment
                                , tags
                            from folder
                            where id_parent = {self.id}
                            """)
                
            self.l_dirs = [Ffolder(*x) for x in self.con.fetchall()]

            cur.execute(f"""
                            select id
                                , id_folder
                                , name
                                , comment
                                , tags
                            from file
                            where id_folder = {self.id}
                            """)
            
            self.l_files = [Ffile(*x) for x in self.con.fetchall()]


    def download(self, down_path: str, todo: bool = False):
        dir_path = down_path + os.sep + self.name
        os.makedirs(dir_path)
        
        for dir in self.l_dirs: dir.download(down_path=dir_path, todo=todo)
        for f in self.l_files: f.download(down_path=dir_path, todo=todo)


    def load(self):
        self.get_elems(self.src_path)


    def save(self):
        cur = self.con.cursor()
        cur.execute("""
                    insert into obj(id_prj
                                        , id_parent
                                        , name
                                        , comment
                                        , object_type
                                        , tags)
                                    values({id_prj}
                                        , {id_parent}
                                        , {name}  
                                        , {comment}
                                        , 0
                                        , {tags})
                        returning id
                    """.format(id_prj = '{id_prj}'.format(id_prj=self.id_prj) if self.id_prj else 'null'
                            , id_parent = '{id_parent}'.format(id_parent=self.id_parent) if self.id_parent else 'null'
                            , name = "'{name}'".format(name=self.name) if self.name else 'null'
                            , comment = "'{comment}'".format(comment=self.comment) if self.comment else 'null'
                            , tags = "'{tags}'".format(tags=self.tags) if self.tags else 'null'))
        
        self.id = cur.fetchone()[0]
        
        for dir in self.l_dirs:
            dir.id_parent = self.id
            dir.id_prj = self.id_prj
            dir.save()

        for f in self.l_files:
            f.id_folder = self.id
            f.id_prj = self.id_prj
            f.save()

        self.con.commit()


class Prj(FileSystem):

    __slots__ = ('id', 'id_item', 'name', 'comment', 'tags')
    
    src_path = None

    l_dirs = list()
    l_files = list()

    def __init__(self, id = None, id_item = None, name = None,  comment = None, tags = None):
        self.id = id
        self.id_item = id_item
        self.name = name
        self.comment = comment
        self.tags = tags

        if self.id:
            cur = self.con.cursor()
            cur.execute(f"""
                            select id
                                , id_prj
                                , id_parent
                                , name
                                , comment
                                , tags
                            from obj
                            where id_prj = '{self.id}'
                              and object_type = 0
                            """)
            
            self.l_dirs = [Ffolder(*x) for x in cur.fetchall()]

            cur.execute(f"""
                            select id
                                , id_prj
                                , id_folder
                                , name
                                , comment
                                , tags
                            from file
                            where id_prj = {self.id}
                              and object_type = 1
                                """)
            
            self.l_files = [Ffile(*x) for x in cur.fetchall()]

    def download(self, down_path: str, todo: bool = False):
        for dir in self.l_dirs: dir.download(down_path=down_path, todo=todo)
        for f in self.l_files: f.download(down_path=down_path, todo=todo)


    def load(self):
        self.get_elems(self.src_path)
        for dir in self.l_dirs:
            dir.src_path = self.src_path + os.sep + dir.name
            dir.load()
        for f in self.l_files:
            f.src_path = self.src_path + os.sep + f.name
            f.load()
    
    def save(self):
        cur = self.con.cursor()
        cur.execute("""
                            insert into prj(id_item
                                                , name
                                                , comment
                                                , tags)
                                            values({id_item}
                                                , {name}  
                                                , {comment}
                                                , {tags})
                                returning id
                            """.format(id_item = '{id_item}'.format(id_item=self.id_item) if self.id_item else 'null'
                                    , name = "'{name}'".format(name=self.name) if self.name else 'null'
                                    , comment = "'{comment}'".format(comment=self.comment) if self.comment else 'null'
                                    , tags = "'{tags}'".format(tags=self.tags) if self.tags else 'null'))

        self.id = cur.fetchone()[0]
        self.con.commit()

        for folder in self.l_dirs:
            folder.id_prj = self.id
            folder.save()

        for file in self.l_files:
            file.id_prj = self.id
            file.save()

class Item(FileSystem):

    __slots__ = ('id', 'name')

    l_prj = list()

    def __init__(self, id, name):
        self.id = id
        self.name = name

        cur = self.con.cursor()
        cur.execute(f"""
                        select id
                            , id_item
                            , name
                            , comment
                            , tags
                        from prj
                        where id_item = {self.id}
                        """)
            
        self.l_prj = [Prj(*x) for x in cur.fetchall()]

    
    def save(self):
        cur = self.con.cursor()
        cur.execute("""
                        insert into item(name
                                            , comment)
                                        values({name}  
                                            , {comment})
                            returning id
                        """.format(name = "'{name}'".format(name=self.name) if self.name else 'null'
                                , comment = "'{comment}'".format(comment=self.comment) if self.comment else 'null'))
            
        self.id = cur.fetchone()[0]
        self.con.commit()

        for prj in self.l_prj:
            prj.id_item = self.id
            prj.save()


def get_file(id = None, name = None, comment = None):
    with connect(db) as con:
        cur = con.cursor()
        cur.execute("""
                    select id
                            , id_folder
                            , name
                            , comment
                            , tags
                        from obj
                        where 1=1 {f_id} {f_name} {f_comment}
                          and object_type = 1
                    """.format(
                        f_id = 'and id = {id}'.format(id = id) if id else ''
                        , f_name = "and lower(name) like lower('%{name}%')".format(name = name) if name else ''
                        , f_comment = "and lower(comment) like lower('%{comment}%')".format(comment = comment) if comment else ''))
    
        return [Ffile(*x) for x in con.fetchall()]
 

def get_elems(obj, src_path):
    path_elems = os.walk(src_path).__next__()
    obj.l_dirs = [Ffolder(name = x) for x in path_elems[1]]
    obj.l_files = [Ffile(name = x) for x in path_elems[2]]
    return obj


def get_item_id(item_name: str, rn: int = None):
    with connect(db) as con:
        cur = con.cursor()
        cur.execute("""
                select id
                   from (select id
                                , row_number() OVER (PARTITION BY (name) ORDER BY NULL) rn
                            from item
                           where lower(name) like lower('%{item_name}%') )
                where 1=1 {f_rn}
                """.format(
                        item_name="{item_name}".format(item_name=item_name)
                        , f_rn="and rn = {rn}".format(rn=rn) if rn else ''))
        return cur.fetchone()[0]


def get_prj(id_prj: int) -> Prj:
    with connect(db) as con:
        cur = con.cursor()
        cur.execute(f"""
                select id
                     , id_item
                     , name
                     , comment
                     , tags
                  from prj
                 where id = {id_prj}
                 """)
        
        return Prj(*cur.fetchone())



def save_prj(src_path, item_name):
    global db
    db = os.path.curdir + os.sep + item_name.upper() + '.db'
    prj_name = os.path.basename(src_path)
    id_item = get_item_id(item_name=item_name, rn = 1)
    prj = Prj(id_item=id_item, name=prj_name)
    prj.src_path = src_path

    prj.load()
    prj.save()


def set_item(item_name: str, comment: str = None):
    with connect(db) as con:
        cur = con.cursor()
        cur.execute("""
                insert into item (name
                                     , comment)
                              values ({item_name}
                                     , {comment})
                """.format(item_name = item_name
                           , comment = "'{comment}'".format(comment=comment) if comment else 'null'))
        con.commit()        
        
def set_comment(type: str, id: int, comment: str) -> None:
    with connect(db) as con:
        cur = con.cursor()
        cur.execute(f"""
                update {type}
                   set comment = '{comment}'
                where id = {id}
                """)
        con.commit()