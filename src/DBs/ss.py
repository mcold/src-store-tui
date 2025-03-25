# coding: utf-8

import oop
import typer
import os

app = typer.Typer()


@app.command(help='Download project')
def dp(down_path: str, id_prj: int, todo: bool = False) -> None:
    prj = oop.get_prj(id_prj=id_prj)
    prj.download(down_path=down_path, todo=todo)


@app.command(help='Items list')
def ils(name: str) -> None:
    # TODO
    pass


@app.command(help='Project list')
def pls(name: str) -> None:
    # TODO
    pass


@app.command(help='Save Project')
def sp(src_path: str, item_name: str) -> None:
    oop.save_prj(src_path=src_path, item_name=item_name)


@app.command(help='Set comment')
def sc(type: str, id: int, comment: str) -> None:
    try:
        match type.lower():
            case 'srcc' | 'folder' | 'file' | 'prj' | 'item':
                oop.set_comment(type, id, comment)
            case _:
                raise ValueError
    except:
        print("Oops!  That was no valid type.  Valid types: srcc, folder, file, prj, item. Try again...")


if __name__ == "__main__":
    app()