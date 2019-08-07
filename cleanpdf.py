#!/usr/local/bin/python
# -*- coding: utf-8 -*-
# Author: James Campbell
# Date Created: May 24th 2016
# Date Updated: February 23rd 2019
# What: Views, Stores, and Wipes Metadata from PDF files

from __future__ import unicode_literals
import codecs
try:
    from tkinter import *
    py3 = True
except:
    from Tkinter import *
    py3 = False
try:
    from tkinter.filedialog import askopenfilename
except:
    from tkFileDialog import askopenfilename
try:
    from tkinter.messagebox import showerror
except:
    from tkMessageBox import showerror
try:
    import tkMessageBox
except:
    import tkinter.messagebox as tkMessageBox
from PyPDF2 import PdfFileReader, PdfFileMerger
try:
    from tkinter import ttk
except:
    import ttk  # not python 3
infos = ''


def center_window(width=300, height=200):
    # get screen width and height
    screen_width = root.winfo_screenwidth()
    screen_height = root.winfo_screenheight()

    # calculate position x and y coordinates
    x = (screen_width / 2) - (width / 2)
    y = (screen_height / 2) - (height / 2)
    root.geometry('%dx%d+%d+%d' % (width, height, x, y))
    a = root.winfo_geometry().split('+')[0]
    b = a.split('x')
    w = int(b[0])
    h = int(b[1])
    root.geometry('%dx%d' % (w+1, h+1))


class MyFrame(Frame):

    def __init__(self, parent):
        Frame.__init__(self)
        #menubar = Menu(self)
        #menubar.add_command(label="Hello!", command=self.load_file)
        self.master.minsize(width=370, height=85)
        self.master.maxsize(width=370, height=85)
        self.master.title("CLEAN PDF")
        self.master.labelText = Text(
            self.master, height=2, width=30)
        self.master.labelText.insert(END, 'Select a PDF to analyze.')
        self.master.rowconfigure(1, weight=1)
        self.master.columnconfigure(4, weight=1)
        self.grid(row=0, column=4)

        self.button = ttk.Button(self, text="ANALYZE MODE - Select PDF to view metadata",
                                 command=self.load_file_pdf, width=36, style="C.TButton")
        self.button.grid(row=0, column=1, sticky="", columnspan=2)
        self.button2 = ttk.Button(self, text="CLEAN MODE - Select PDF to wipe metadata",
                                  command=self.wipe_file_pdf, width=36, style="C.TButton")
        self.button2.grid(row=1, column=1, sticky="", columnspan=2)

    def load_file_pdf(self):
        self.button.configure()
        fname = askopenfilename(filetypes=(("PDF files", "*.pdf"),
                                           ("All files", "*.*")))
        if fname:
            #print("""here it comes: self.settings["template"].set(fname)""")
            # print(fname)  - testing only -jc
            oldname = fname.rsplit('/', 1)[1]
            pather = fname.rsplit('/', 1)[0]
            cleanname = oldname.rsplit('.', 1)[0] + '-clean.pdf'
            metatxt = oldname.rsplit('.', 1)[0] + '-metadata.txt'
            try:
                reader = PdfFileReader(open(fname, "rb"))
            except UnicodeEncodeError as e:
                tkMessageBox.showinfo(
                    "Error", "The chars in this filename broke this app:\n\n{}".format(oldname))
                return
            # get the metadata available for the pdf - jc
            try:
                infos = reader.getDocumentInfo()
            except:
                tkMessageBox.showinfo(
                    "Done", "No metadata found for:\n\n{}".format(oldname))
                return
            # infos is now a dictionary
            cur_metadata_str = ''  # a way to show all the current metadata fields -jc
            cleandict = dict()
            olddict = infos
            for value in infos:
                # print (value, infos[value])  - testing only -jc
                newdictvalue = {value: ''}
                cleandict.update(newdictvalue)
                valued = ''.join(value.split('/', 1)
                                 [1]).encode('utf8', 'ignore')
                try:
                    infoed = infos[value].encode('utf8', 'ignore')
                except:
                    infoed = infos[value]
                try:
                    cur_metadata_str = cur_metadata_str + \
                        valued.encode('utf8', 'ignore') + ' : ' + infoed + '\n'
                except TypeError as e:  # handle weird ArrayObject error
                    infoed = ''.join(infos[value]).encode('utf8', 'ignore')
                    cur_metadata_str = cur_metadata_str + \
                        valued.encode('utf8', 'ignore') + ':' + infoed + '\n'
                except UnicodeDecodeError as e:  # handle weird ArrayObject error
                    infoed = ''.join(infos[value])
                    cur_metadata_str = cur_metadata_str + valued + ' : ' + infoed + '\n'
            if tkMessageBox.askyesno("Proceed", "Found {} items in metadata. \n{}\nProceed to store metadata?".format(len(infos), cur_metadata_str, cleanname)):
                with codecs.open(pather + '/' + metatxt, 'w', encoding='utf8') as metasavefile:
                    for v in olddict:
                        try:
                            vchanged = v.split(
                                '/', 1)[1].encode('utf8', 'ignore')
                        except:
                            vchanged = v.split('/', 1)[1]
                        try:
                            metasavefile.write(vchanged.encode(
                                'utf8', 'ignore') + ',' + olddict[v].encode('utf8', 'ignore') + '\n')
                        except TypeError as e:
                            infoed = ''.join(olddict[v]).encode(
                                'utf8', 'ignore')
                            metasavefile.write(vchanged.encode(
                                'utf8', 'ignore') + ',' + infoed + '\n')
                        except AttributeError as e:
                            infoed = ''.join(olddict[v]).encode(
                                'utf8', 'ignore')
                            metasavefile.write(vchanged.encode(
                                'utf8', 'ignore') + ',' + infoed + '\n')
                        except UnicodeDecodeError as e:
                            infoed = ''.join(olddict[v])
                            metasavefile.write(vchanged + ',' + infoed + '\n')
                tkMessageBox.showinfo("Complete", "Metadata saved as\n{}".format(
                    pather + '/' + metatxt, 'w'))
                return
            else:
                tkMessageBox.showinfo(
                    "Canceled", "Metadata shown but not saved.")
                return

    def wipe_file_pdf(self):
        self.button.configure()
        fname = askopenfilename(filetypes=(("PDF files", "*.pdf"),
                                           ("All files", "*.*")))
        if fname:
            #print("""here it comes: self.settings["template"].set(fname)""")
            # print(fname)  - testing only -jc
            oldname = fname.rsplit('/', 1)[1]
            pather = fname.rsplit('/', 1)[0]
            cleanname = oldname.rsplit('.', 1)[0] + '-clean.pdf'
            metatxt = oldname.rsplit('.', 1)[0] + '-metadata.txt'
            try:
                reader = PdfFileReader(open(fname, "rb"))
            except UnicodeEncodeError as e:
                tkMessageBox.showinfo(
                    "Error", "The chars in this filename broke this app:\n\n{}".format(oldname))
                return
            # get the metadata available for the pdf - jc
            try:
                infos = reader.getDocumentInfo()
            except:
                tkMessageBox.showinfo(
                    "Done", "No metadata found for:\n\n{}".format(oldname))
                return
            # infos is now a dictionary
            cur_metadata_str = ''  # a way to show all the current metadata fields -jc
            cleandict = dict()
            olddict = infos
            for value in infos:
                # print (value, infos[value])  - testing only -jc
                newdictvalue = {value: ''}
                cleandict.update(newdictvalue)
                valued = ''.join(value.split('/', 1)
                                 [1]).encode('utf8', 'ignore')
                try:
                    infoed = infos[value].encode('utf8', 'ignore')
                except:
                    infoed = infos[value]
                try:
                    cur_metadata_str = cur_metadata_str + \
                        valued.encode('utf8', 'ignore') + ' : ' + infoed + '\n'
                except TypeError as e:  # handle weird ArrayObject error
                    infoed = ''.join(infos[value]).encode('utf8', 'ignore')
                    cur_metadata_str = cur_metadata_str + \
                        valued.encode('utf8', 'ignore') + ':' + infoed + '\n'
                except UnicodeDecodeError as e:  # handle weird ArrayObject error
                    infoed = ''.join(infos[value])
                    cur_metadata_str = cur_metadata_str + valued + ' : ' + infoed + '\n'
            if tkMessageBox.askyesno("Proceed", "Found {} items in metadata. \n{}\nProceed to wipe metadata?\n File will be saved as:\n\n{}".format(len(infos), cur_metadata_str, cleanname)):

                # create a new merge object and append
                # the opened data to it with updated blank infos -jc
                writr = PdfFileMerger()
                writr.append(reader)
                if len(infos) > 0:
                    infos = cleandict
                else:
                    infos = {u'/Trapped': '', u'/Title': '', u'/Author': '', u'/Subject': '',
                             u'/Producer': '', u'/Content creator': '', u'/CreationDate': '', u'/ModDate': ''}

                writr.addMetadata(infos)
                writr.write(pather + '/' + cleanname)
                tkMessageBox.showinfo(
                    "Complete", "File cleaned and saved as:\n\n{}".format(cleanname))

            else:
                tkMessageBox.showinfo(
                    "Canceled", "Operation canceled. No clean file created.")


if __name__ == "__main__":
    root = Tk()
    root.lift()
    style = ttk.Style()
    style.map("C.TButton",
              foreground=[('pressed', 'red'), ('active', 'blue')],
              background=[('pressed', '!disabled', 'black'),
                          ('active', 'white')]
              )
    root.call('wm', 'attributes', '.', '-topmost', True)
    root.after_idle(root.call, 'wm', 'attributes', '.', '-topmost', False)
    center_window()
    app = MyFrame(root)
    app.mainloop()
